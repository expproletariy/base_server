package chat

import (
	"github.com/expproletariy/base_server/types"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

//Room - type which represent group of clients
type Room struct {
	types.Room
	//Group of clients
	clients map[string]*Client

	mutex sync.RWMutex
}

//NewRoom - create new Room
func NewRoom(room types.Room) *Room {
	return &Room{
		Room:    room,
		clients: make(map[string]*Client),
	}
}

//Register new client in the room
func (room *Room) Register(client *Client) {
	room.mutex.Lock()
	defer room.mutex.Unlock()
	//Disconnect client if no activity for 10 minutes
	client.conn.SetReadDeadline(time.Now().Add(time.Minute * 10))

	room.clients[client.ID] = client
}

func (room *Room) WatchClientMessages(userID string, save func(message types.Message) error) error {
	room.mutex.RLock()
	client, ok := room.clients[userID]
	room.mutex.RUnlock()
	closeSignal := make(chan byte)
	pingTicker := time.NewTicker(time.Second * 10)
	defer pingTicker.Stop()
	go func() {
		for range pingTicker.C {
			err := client.conn.WriteControl(websocket.PingMessage, []byte(""), time.Now().Add(time.Second*10))
			if err != nil {
				closeSignal <- 1
				return
			}
		}
	}()
	if ok {
		for {
			select {
			case <-closeSignal:
				room.Remove(client.ID)
				return nil
			default:
				msg, err := client.ReadMessage()
				if err != nil {
					room.Remove(client.ID)
					return err
				}
				msg.CreatedAt = time.Now()
				msg.UserID = client.ID
				msg.RoomID = room.ID
				msg.ID = uuid.NewV4().String()
				msg.UserName = client.Name
				err = save(msg)
				if err != nil {
					room.Remove(client.ID)
					return err
				}
				err = room.Message(msg)
				if err != nil {
					room.Remove(client.ID)
					return err
				}

			}
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
func (room *Room) Message(message types.Message) error {
	room.mutex.RLock()
	defer room.mutex.RUnlock()

	for _, client := range room.clients {
		//if message.UserID != client.ID {
		//
		//}
		err := client.conn.WriteJSON(message)
		if err != nil {
			return err
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

//CheckUser for active connection in current room
func (room *Room) CheckUser(id string) bool {
	room.mutex.RLock()
	defer room.mutex.RUnlock()
	_, ok := room.clients[id]
	return ok
}
