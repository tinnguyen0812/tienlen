package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"tienlen-server/managers"
	"tienlen-server/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade"})
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
		ConnID:   conn.RemoteAddr().String(),
		Username: username,
		Conn:     conn,
	}
	room.AddPlayer(player)

	room.Broadcast("player_joined", gin.H{
		"id":       playerID,
		"username": username,
	}, nil)

	// Khởi tạo ván chơi nếu đủ người
	if len(room.Players) >= 2 {
		startGame(room)
	}

	// Lắng nghe message từ client
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}

		switch msg.Type {
		case "play_cards":
			handlePlayCards(room, player, msg.Data)
		}
	}
}

func startGame(room *models.Room) {
	deck := models.Shuffle(models.GenerateDeck())
	playerIDs := room.GetPlayerIDs()
	cardsPerPlayer := len(deck) / len(playerIDs)

	for i, id := range playerIDs {
		hand := deck[i*cardsPerPlayer : (i+1)*cardsPerPlayer]
		room.PlayerCards[id] = hand

		player := room.Players[id]
		player.Conn.WriteJSON(gin.H{
			"type": "your_cards",
			"data": hand,
		})
	}

	// Bắt đầu lượt chơi
	room.TurnManager.Start(playerIDs)
	first := room.TurnManager.GetCurrentPlayer()
	room.Broadcast("start_turn", gin.H{"playerId": first}, nil)
}

func handlePlayCards(room *models.Room, player *models.Player, rawData json.RawMessage) {
	var playedCards []models.Card
	_ = json.Unmarshal(rawData, &playedCards)

	current := room.TurnManager.GetCurrentPlayer()
	if player.ID != current {
		player.Conn.WriteJSON(gin.H{"error": "Không phải lượt của bạn"})
		return
	}

	if !IsValidPlay(playedCards, room.LastPlay) {
		player.Conn.WriteJSON(gin.H{"error": "Bài không hợp lệ"})
		return
	}

	room.LastPlay = playedCards
	room.PlayerCards[player.ID] = removeCards(room.PlayerCards[player.ID], playedCards)

	room.Broadcast("player_played", gin.H{
		"playerId": player.ID,
		"cards":    playedCards,
	}, nil)

	if len(room.PlayerCards[player.ID]) == 0 {
		room.Broadcast("player_won", gin.H{"playerId": player.ID}, nil)
		// TODO: Kiểm tra xem có kết thúc ván chơi không
		return
	}

	room.TurnManager.NextTurn()
	next := room.TurnManager.GetCurrentPlayer()
	room.Broadcast("next_turn", gin.H{"playerId": next}, nil)
}
