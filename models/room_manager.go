package models

import (
	"sync"
)

type RoomManager struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

var roomManager *RoomManager

func init() {
	roomManager = &RoomManager{
		rooms: make(map[string]*Room),
	}
}

func GetRoomManager() *RoomManager {
	return roomManager
}

func (m *RoomManager) GetRoom(id string) (*Room, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	room, exists := m.rooms[id]
	return room, exists
}

func (m *RoomManager) CreateRoom(id string) *Room {
	m.mu.Lock()
	defer m.mu.Unlock()
	room := NewRoom(id)
	m.rooms[id] = room
	return room
}
