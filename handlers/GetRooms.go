package handlers

import (
	"github.com/expproletariy/base_server/chat"
	"github.com/labstack/echo"
	"net/http"
)

type Room struct {
	//Room id
	ID string `json:"id"`

	//Room name
	Name string `json:"name"`
}

func GetRooms(c echo.Context) error {

	rooms := make([]Room, 0, 10)

	chat.Manager().EachRoom(func(id, name string) {
		rooms = append(rooms, Room{ID: id, Name: name})
	})

	return c.JSON(http.StatusOK, map[string][]Room{
		"rooms": rooms,
	})

}
