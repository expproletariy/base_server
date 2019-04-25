package handlers

import (
	"encoding/json"
	"github.com/expproletariy/base_server/models"
	"github.com/expproletariy/base_server/session"
	"github.com/expproletariy/base_server/types"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"net/http"
	"strings"
	"time"
)

func CreateRoom(c echo.Context) error {

	body := models.Room{}

	err := json.NewDecoder(c.Request().Body).Decode(&body)

	if err != nil {
		return err
	}

	if len(body.Name) == 0 {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			`Need to set "name" of room param`,
		)
	}
	_, err = models.StmtCheckRoomName.Exec(body.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, "Room name already exists")
	}

	AuthHeader := c.Request().Header.Get("Authorization")
	JWT := strings.Split(AuthHeader, " ")

	userInfo, ok := session.GetClaimsInfo(JWT[1])
	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, "Invalid token or invalid user!")
	}

	dbCtx, ok := models.GetContext()
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Broken db connection")
	}

	tx, err := dbCtx.Beginx()
	if err != nil {
		return err
	}

	newRoom := models.Room{Room: types.Room{
		Name: body.Name,
	}}

	newRoom.ID = uuid.NewV4().String()
	newRoom.CreatedAt = time.Now().UTC()

	stmtNewRoom, err := tx.PrepareNamed("INSERT INTO rooms (id, name, created_at) VALUES (:id, :name, :created_at)")
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = stmtNewRoom.Query(newRoom)
	if err != nil {
		tx.Rollback()
		return err
	}

	newUserRoom := models.UserRoom{
		UserID: userInfo.UserID,
		RoomID: newRoom.ID,
		User: types.User{
			Creator: true,
		},
	}

	stmtNewUserRoom, err := tx.PrepareNamed("INSERT INTO user_room " +
		"(users_id, rooms_id, creator) VALUES (:users_id, :rooms_id, :creator)",
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmtNewUserRoom.Query(newUserRoom)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"ok":      true,
		"room_id": newRoom.ID,
	})
}
