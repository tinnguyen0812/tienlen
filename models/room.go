package models

import (
	"github.com/gin-gonic/gin"
)

type Room struct {
	ID          string
	Players     map[string]*Player
	PlayerCards map[string][]Card
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
