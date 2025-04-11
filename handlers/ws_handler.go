package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"tienlen-server/managers"
	"tienlen-server/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "WebSocket upgrade failed"})
		return
	}

	playerID := c.Query("player_id")
	username := c.Query("username")
	roomID := c.Query("room_id")

	roomManager := managers.GetRoomManager()
	room, exists := roomManager.GetRoom(roomID)
	if !exists {
		room = roomManager.CreateRoom(roomID)
	}
	player := &models.Player{
		ID:       playerID,
		Username: username,
		Conn:     conn,
	}

	room.AddPlayer(player)

	// Gửi danh sách người chơi cho các client
	room.Broadcast("player_joined", gin.H{
		"id":       playerID,
		"username": username,
	}, nil)

	// Khi đủ 2 người trở lên, chia bài
	if len(room.Players) >= 2 {
		deck := models.Shuffle(models.GenerateDeck())
		i := 0
		for _, p := range room.Players {
			hand := deck[i*13 : (i+1)*13]
			room.PlayerCards[p.ID] = hand
			p.Conn.WriteJSON(gin.H{"type": "your_cards", "data": hand})
			i++
		}
		room.DetermineFirstPlayer()

		room.Broadcast("game_start", gin.H{
			"first_player_id":   room.CurrentTurnPlayerID,
			"first_player_name": room.Players[room.CurrentTurnPlayerID].Username,
		}, nil)

	}
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read error: %v", err)
			room.RemovePlayer(player.ID)
			break
		}

		var payload map[string]interface{}
		if err := json.Unmarshal(msg, &payload); err != nil {
			log.Println("invalid message format:", err)
			continue
		}

		eventType, _ := payload["type"].(string)
		data := payload["data"]

		switch eventType {
		case "play_card":
			var cards []string
			if arr, ok := data.([]interface{}); ok {
				for _, c := range arr {
					if s, ok := c.(string); ok {
						cards = append(cards, s)
					}
				}
			}
			room.HandlePlayCard(player.ID, cards)

			// 👉 Bạn có thể thêm các case như "pass", "chat", v.v. ở đây
		}
	}
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
