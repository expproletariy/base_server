package chat

import (
	"github.com/satori/go.uuid"
	"sync"
)

//Manager of chat, should be one for app
type manager struct {
	//Rooms of chat
	rooms map[string]*Room

	mutex sync.RWMutex
}

//App instance of manager
var instance *manager

var once = sync.Once{}

//Manager - function to get instance of manager
func Manager() *manager {
	once.Do(func() {
		instance = &manager{
			rooms: make(map[string]*Room),
		}
	})
	return instance
}

//CreateRoom init new chat room inside the manager
func (mng *manager) CreateRoom(name string) (string, error) {
	mng.mutex.Lock()
	defer mng.mutex.Unlock()
	roomID := uuid.NewV4().String()
	mng.rooms[roomID] = NewRoom(roomID, name)

	return roomID, nil
}

//DeleteRoom remove room and rooms clients (close socket connections)
func (mng *manager) DeleteRoom(id string) error {
	mng.mutex.Lock()
	defer mng.mutex.Unlock()
	if room, ok := mng.rooms[id]; ok {
		err := room.Clear()
		if err != nil {
			return err
		}
		delete(mng.rooms, room.ID)
	}
	return NewError("Try to delete nonexistent room")
}

//GetRoom returns existing room by id
func (mng *manager) GetRoom(id string) (*Room, error) {
	mng.mutex.RLock()
	defer mng.mutex.RUnlock()
	if room, ok := mng.rooms[id]; ok {
		return room, nil
	}
	return nil, NewError("Try to get nonexistent room")
}

//GetRoom returns existing rooms
func (mng *manager) EachRoom(f func(id, name string)) {
	for _, room := range mng.rooms {
		mng.mutex.RLock()
		f(room.ID, room.Name)
		mng.mutex.RUnlock()
	}
}
