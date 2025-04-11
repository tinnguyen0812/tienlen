package managers

import (
	"sync"
	"tienlen-server/models"
)

type RoomManager struct {
	Rooms map[string]*models.Room
	Mutex sync.Mutex
}

var instance *RoomManager
var once sync.Once

func GetRoomManager() *RoomManager {
	once.Do(func() {
		instance = &RoomManager{
			Rooms: make(map[string]*models.Room),
		}
	})
	return instance
}

func (rm *RoomManager) CreateRoom(roomID string) *models.Room {
	rm.Mutex.Lock()
	defer rm.Mutex.Unlock()
	room := models.NewRoom(roomID)
	rm.Rooms[roomID] = room
	return room
}

func (rm *RoomManager) GetRoom(roomID string) (*models.Room, bool) {
	rm.Mutex.Lock()
	defer rm.Mutex.Unlock()
	room, exists := rm.Rooms[roomID]
	return room, exists
}
