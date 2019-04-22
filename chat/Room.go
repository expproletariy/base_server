package chat

import (
	"fmt"
	"sync"
	"time"
)

//Room - type which represent group of clients
type Room struct {
	//Group of clients
	clients map[string]*Client

	//Room id
	ID string

	//Room name
	Name string

	mutex sync.RWMutex
}

//NewRoom - create new Room
func NewRoom(id, name string) *Room {
	return &Room{
		Name:    name,
		ID:      id,
		clients: make(map[string]*Client),
	}
}

//Register new client in the room
func (room *Room) Register(client *Client) {
	room.mutex.Lock()
	defer room.mutex.Unlock()

	room.clients[client.ID] = client
}

func (room *Room) WatchClientMessages(id string) error {
	room.mutex.RLock()
	client, ok := room.clients[id]
	room.mutex.RUnlock()
	if ok {
		for {
			msg, err := client.ReadMessage()
			if err != nil {
				room.Remove(client.ID)
				return err
			}
			msg.Time = time.Now()
			room.Message(msg)
		}
	}
	return nil
}

//Remove client from room
func (room *Room) Remove(id string) error {
	room.mutex.Lock()
	defer room.mutex.Unlock()

	if client, ok := room.clients[id]; ok {
		err := client.Close()
		if err != nil {
			return nil
		}
		delete(room.clients, client.ID)
	}

	return nil
}

//Message to the room clients
func (room *Room) Message(message Message) error {
	room.mutex.RLock()
	defer room.mutex.RUnlock()

	for _, client := range room.clients {
		fmt.Println(message)
		if message.UserID != client.ID {
			err := client.conn.WriteJSON(message)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//Clear removes all clients and close connections
func (room *Room) Clear() error {
	room.mutex.Lock()
	defer room.mutex.Unlock()
	for _, client := range room.clients {
		err := client.Close()
		if err != nil {
			return nil
		}
		delete(room.clients, client.ID)
	}
	return nil
}
