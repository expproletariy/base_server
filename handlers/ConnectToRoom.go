package handlers

import (
	"github.com/expproletariy/base_server/chat"
	"github.com/expproletariy/base_server/models"
	"github.com/expproletariy/base_server/session"
	"github.com/expproletariy/base_server/types"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
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
	//Upgrade connection to set up websocket mod
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	//Get from params room_id and user-id
	roomID := c.QueryParam("room_id")
	userInfo, ok := session.GetClaimsInfo(c.QueryParam("token"))

	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, "Invalid token or invalid user!")
	}

	//Check room in cache if not try to search in db or fail
	room, err := chat.Manager().GetRoom(roomID)
	if err != nil {
		var dbRoom models.Room
		err = models.StmtGetRoom.Get(&dbRoom, roomID)
		if err != nil {
			return err
		}
		_, err = chat.Manager().RegisterRoom(dbRoom.Room)
		if err != nil {
			return err
		}
		room, err = chat.Manager().GetRoom(roomID)
		if err != nil {
			return err
		}
	}

	//Check user, if user already connected fall down
	if room.CheckUser(userInfo.UserID) {
		return echo.NewHTTPError(http.StatusForbidden, "You was trying to connect to the room where you are!")
	}

	//Try to get user info from db to create new connection, if not fall down
	var user models.User
	err = models.StmtGetUserRoom.Get(&user, userInfo.UserID, roomID)
	if err != nil {
		_, err = models.StmtNewUserRoomByDefault.Exec(models.UserRoom{
			RoomID: roomID,
			UserID: userInfo.UserID,
		})
		if err != nil {
			return err
		}
		err = models.StmtGetUserRoom.Get(&user, userInfo.UserID, roomID)
		if err != nil {
			return err
		}
	}

	if user.Blocked {
		return echo.NewHTTPError(http.StatusForbidden, "You was trying to connect to the room where you was blocked!")
	}

	room.Register(chat.NewClient(ws, user.User))

	c.Logger().Printf("user %s connected to the room %s\n", userInfo.UserID, roomID)

	msg := models.Message{
		Message: types.Message{
			SimpleMessage: types.SimpleMessage{
				CreatedAt: time.Now().UTC(),
				ID:        uuid.NewV4().String(),
				Text:      "connected",
				UserID:    userInfo.UserID,
			},
			RoomID:   roomID,
			UserName: user.Name,
		},
	}

	_, err = models.StmtSaveMessage.Exec(msg)
	if err != nil {
		return err
	}
	err = room.Message(msg.Message)
	if err != nil {
		return err
	}
	err = room.WatchClientMessages(userInfo.UserID, func(message types.Message) error {
		_, err = models.StmtSaveMessage.Exec(message)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
