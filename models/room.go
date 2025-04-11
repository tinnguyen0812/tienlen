package models

import (
	"sync"
)

type Room struct {
	ID          string
	Players     map[string]*Player
	PlayerCards map[string][]Card
	TurnManager *TurnManager
	LastPlay    []Card
	Mutex       sync.Mutex
}

func NewRoom(id string) *Room {
	return &Room{
		ID:      id,
		Players: make(map[string]*Player),
	}
}

func (r *Room) AddPlayer(player *Player) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.Players[player.ID] = player
}

func (r *Room) RemovePlayer(playerID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	delete(r.Players, playerID)
}

func (r *Room) Broadcast(messageType string, data interface{}, sendFunc func(player *Player, messageType string, data interface{})) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	for _, player := range r.Players {
		sendFunc(player, messageType, data)
	}
}
