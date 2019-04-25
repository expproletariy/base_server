package handlers

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/expproletariy/base_server/models"
	"github.com/expproletariy/base_server/session"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

func SignIn(c echo.Context) error {

	body := models.User{}

	err := json.NewDecoder(c.Request().Body).Decode(&body)

	if len(body.Password) == 0 {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			`Need to set "password" param`,
		)
	}
	if len(body.Name) == 0 {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			`Need to set "name" param`,
		)
	}

	body.ID = uuid.NewV4().String()

	hash := sha512.New()
	hash.Write([]byte(body.Password))
	body.Password = fmt.Sprintf("%x", hash.Sum(nil))

	_, err = models.StmtNewUser.Exec(&body)
	if err != nil {
		return err
	}

	token, err := session.NewToken(session.JWTClaims{
		UserID:   body.ID,
		UserName: body.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 168).Unix(),
		},
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token":   token,
		"user_id": body.ID,
	})
}
