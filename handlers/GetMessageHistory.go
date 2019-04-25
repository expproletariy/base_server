package handlers

import (
	"github.com/expproletariy/base_server/models"
	"github.com/expproletariy/base_server/session"
	"github.com/expproletariy/base_server/types"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"strings"
)

func GetMessageHistory(c echo.Context) error {

	queryPage := c.QueryParam("page")
	roomID := c.QueryParam("room_id")
	page := 0
	if len(queryPage) != 0 {
		page, _ = strconv.Atoi(queryPage)
	}
	if page != 0 {
		page -= 1
	}
	page *= 20

	AuthHeader := c.Request().Header.Get("Authorization")
	JWT := strings.Split(AuthHeader, " ")

	userInfo, ok := session.GetClaimsInfo(JWT[1])
	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, "Invalid token or invalid user!")
	}

	userRoom := models.UserRoom{}

	err := models.StmtGetUserRoom.Get(&userRoom, userInfo.UserID, roomID)
	if err != nil {
		return err
	}

	var messages []models.Message
	if userRoom.Blocked {
		err := models.StmtGetMessageHistoryForBlocked.Select(&messages, userInfo.UserID, roomID, page)
		if err != nil {
			return err
		}
	} else {
		err := models.StmtGetMessageHistory.Select(&messages, roomID, page)
		if err != nil {
			return err
		}
	}

	var responseMessages []types.SimpleMessage

	if len(messages) != 0 {
		for _, msg := range messages {
			responseMessages = append(responseMessages, msg.SimpleMessage)
		}
	}

	return c.JSON(http.StatusOK, map[string][]types.SimpleMessage{
		"messages": responseMessages,
	})

}
