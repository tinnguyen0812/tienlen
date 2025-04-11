package models

import (
	"github.com/gin-gonic/gin"
	"sort"
	"tienlen-server/utils"
	"time"
)

type Room struct {
	ID                  string
	Players             map[string]*Player
	PlayerCards         map[string][]Card
	CurrentTurnPlayerID string
	LastPlayedCards     []string // các lá bài vừa được đánh
	LastPlayedPlayerID  string   // người vừa đánh bài
	Winners             []string
}

func NewRoom(id string) *Room {
	return &Room{
		ID:          id,
		Players:     make(map[string]*Player),
		PlayerCards: make(map[string][]Card),
	}
}

func (r *Room) AddPlayer(p *Player) {
	r.Players[p.ID] = p
}

func (r *Room) Broadcast(eventType string, data interface{}, exclude func(*Player, string, interface{}) bool) {
	for _, p := range r.Players {
		if exclude != nil && exclude(p, eventType, data) {
			continue
		}
		p.Conn.WriteJSON(gin.H{
			"type": eventType,
			"data": data,
		})
	}
}

func (r *Room) DetermineFirstPlayer() {
	threeOfSpades := Card{Suit: "Spades", Rank: "3"}

	for playerID, cards := range r.PlayerCards {
		for _, card := range cards {
			if card == threeOfSpades {
				r.CurrentTurnPlayerID = playerID
				return
			}
		}
	}
}
func (r *Room) HandlePlayCard(playerID string, cards []string) {
	// 1. Kiểm tra đúng lượt
	if r.CurrentTurnPlayerID != playerID {
		r.Players[playerID].Conn.WriteJSON(gin.H{
			"type":  "error",
			"error": "Không phải lượt của bạn",
		})
		return
	}
	if len(cards) == 0 {
		// Người đánh cuối không được pass
		if r.LastPlayedPlayerID == playerID || len(r.LastPlayedCards) == 0 {
			r.Players[playerID].Conn.WriteJSON(gin.H{
				"type":  "error",
				"error": "Bạn không thể bỏ lượt khi chính bạn là người đánh cuối",
			})
			return
		}

		// Gửi thông báo bỏ lượt
		r.Broadcast("player_passed", gin.H{
			"player_id": playerID,
		}, nil)

		// Chuyển lượt
		r.MoveToNextPlayer()

		// Nếu người chơi tiếp theo là người đánh cuối cùng → reset vòng
		if r.CurrentTurnPlayerID == r.LastPlayedPlayerID {
			r.LastPlayedCards = []string{}
			r.LastPlayedPlayerID = ""
			r.Broadcast("new_round", gin.H{
				"message": "Vòng mới bắt đầu",
			}, nil)
		}

		// Thông báo lượt tiếp theo
		r.Broadcast("next_turn", gin.H{
			"player_id": r.CurrentTurnPlayerID,
		}, nil)

		return
	}
	// 2. Kiểm tra hợp lệ bài
	if !utils.IsValidPlay(cards, r.LastPlayedCards) {
		r.Players[playerID].Conn.WriteJSON(gin.H{
			"type":  "error",
			"error": "Bài không hợp lệ",
		})
		return
	}

	// 3. Cập nhật bài đã đánh cuối cùng
	r.LastPlayedCards = cards
	r.LastPlayedPlayerID = playerID

	// 4. Gửi thông báo đến tất cả người chơi
	r.Broadcast("play_card", gin.H{
		"player_id": playerID,
		"cards":     cards,
	}, nil)

	// 5. Xóa bài đã đánh khỏi bộ bài người chơi
	currentCards := r.PlayerCards[playerID]
	newCards := removeCards(currentCards, cards)
	r.PlayerCards[playerID] = newCards

	// 6. Kiểm tra nếu người chơi hết bài → thắng
	if len(r.PlayerCards[playerID]) == 0 {
		r.Winners = append(r.Winners, playerID)
		r.Broadcast("player_won", gin.H{
			"player_id":   playerID,
			"player_name": r.Players[playerID].Username,
		}, nil)

		// Nếu chỉ còn một người chưa thắng → kết thúc ván
		if len(r.Winners) == len(r.Players)-1 {
			var losers []gin.H
			for _, p := range r.Players {
				if !contains(r.Winners, p.ID) {
					losers = append(losers, gin.H{
						"player_id":   p.ID,
						"player_name": p.Username,
						"cards_left":  r.PlayerCards[p.ID],
					})
				}
			}
			r.Broadcast("game_over", gin.H{
				"winners": r.Winners,
				"losers":  losers,
			}, nil)

			// Reset ván chơi sau vài giây
			go func() {
				time.Sleep(5 * time.Second)
				r.ResetGame()
			}()
			return
		}
	}

	// 7. Chuyển lượt cho người tiếp theo
	r.MoveToNextPlayer()

	// 8. Thông báo lượt tiếp theo
	r.Broadcast("next_turn", gin.H{
		"player_id": r.CurrentTurnPlayerID,
	}, nil)
}

func (r *Room) MoveToNextPlayer() {
	playerIDs := make([]string, 0, len(r.Players))
	for id := range r.Players {
		playerIDs = append(playerIDs, id)
	}

	sort.Strings(playerIDs) // đảm bảo thứ tự nhất quán
	for i, id := range playerIDs {
		if id == r.CurrentTurnPlayerID {
			r.CurrentTurnPlayerID = playerIDs[(i+1)%len(playerIDs)]
			break
		}
	}
}
func (r *Room) RemovePlayer(playerID string) {
	delete(r.Players, playerID)
	delete(r.PlayerCards, playerID)

	// Nếu người chơi bị xóa là người đang đến lượt → chuyển lượt
	if r.CurrentTurnPlayerID == playerID {
		r.MoveToNextPlayer()
	}

	// Thông báo cho các player còn lại
	r.Broadcast("player_left", gin.H{
		"id": playerID,
	}, nil)
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func (r *Room) ResetGame() {
	r.PlayerCards = make(map[string][]Card)
	r.CurrentTurnPlayerID = ""
	r.LastPlayedCards = []string{}
	r.LastPlayedPlayerID = ""
	r.Winners = []string{}
}

func removeCards(hand []Card, toRemove []string) []Card {
	var result []Card
	removeMap := map[string]bool{}
	for _, c := range toRemove {
		removeMap[c] = true
	}

	for _, c := range hand {
		cardStr := c.Rank + c.Suit
		if !removeMap[cardStr] {
			result = append(result, c)
		}
	}
	return result
}
