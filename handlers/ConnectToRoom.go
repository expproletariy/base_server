package handlers

import (
	"fmt"
	"github.com/expproletariy/base_server/chat"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func ConnectToRoom(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	//defer ws.Close()

	roomID := c.QueryParam("room_id")
	room, err := chat.Manager().GetRoom(roomID)
	if err != nil {
		return err
	}
	client := chat.NewClient(ws, uuid.NewV4().String(), "tester")
	room.Register(client)
	fmt.Printf("create client %s in room %s\n", client.ID, roomID)
	msg := chat.Message{
		UserID:   client.ID,
		UserName: client.Name,
		Text:     "connected",
		Time:     time.Now(),
	}
	ws.WriteJSON(msg)
	room.Message(msg)
	err = room.WatchClientMessages(client.ID)
	return err
}
