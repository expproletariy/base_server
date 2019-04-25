package handlers

import (
	"github.com/expproletariy/base_server/models"
	"github.com/labstack/echo"
	"net/http"
)

func GetRooms(c echo.Context) error {

	var rooms []models.Room
	err := models.StmtGetRooms.Select(&rooms)
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, map[string][]models.Room{
		"rooms": rooms,
	})

}
