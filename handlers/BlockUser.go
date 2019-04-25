package handlers

import (
	"encoding/json"
	"github.com/expproletariy/base_server/chat"
	"github.com/expproletariy/base_server/models"
	"github.com/expproletariy/base_server/session"
	"github.com/labstack/echo"
	"net/http"
	"strings"
	"time"
)

func BlockUser(c echo.Context) error {

	body := models.UserRoom{}

	err := json.NewDecoder(c.Request().Body).Decode(&body)

	if err != nil {
		return err
	}

	if len(body.RoomID) == 0 {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			`Need to set room "id" param`,
		)
	}

	if len(body.UserID) == 0 {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			`Need to set blocked user "id" param`,
		)
	}

	AuthHeader := c.Request().Header.Get("Authorization")
	JWT := strings.Split(AuthHeader, " ")

	userInfo, ok := session.GetClaimsInfo(JWT[1])
	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, "Invalid token or invalid user!")
	}

	userRoom := models.UserRoom{}
	err = models.StmtGetUserRoom.Get(&userRoom, userInfo.UserID, body.RoomID)
	if err != nil {
		return err
	}

	if !userRoom.Creator {
		return echo.NewHTTPError(http.StatusForbidden, "To block user you need to be a creator!")
	}

	if userInfo.UserID == body.UserID {
		return echo.NewHTTPError(http.StatusForbidden, "You can not block yourself!")
	}

	room, err := chat.Manager().GetRoom(body.RoomID)
	if err == nil {
		err = room.Remove(body.UserID)
		if err != nil {
			return err
		}
	}

	_, err = models.StmtUserBlocker.Exec(body.Blocked, time.Now().UTC(), body.UserID, body.RoomID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]bool{
		"ok": true,
	})
}
