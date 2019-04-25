package handlers

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/expproletariy/base_server/models"
	"github.com/expproletariy/base_server/session"
	"github.com/labstack/echo"
	"net/http"
	"strings"
	"time"
)

func Login(c echo.Context) error {

	body := models.User{}

	err := json.NewDecoder(c.Request().Body).Decode(&body)

	if len(body.Name) == 0 {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			`Need to set "username" param`,
		)
	}
	if len(body.Password) == 0 {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			`Need to set "password" param`,
		)
	}
	user := models.User{}

	err = models.StmtGetUserByName.Get(&user, body.Name)

	if err != nil {
		return err
	}

	hash := sha512.New()
	hash.Write([]byte(body.Password))
	body.Password = fmt.Sprintf("%x", hash.Sum(nil))

	if strings.ToLower(user.Password) != strings.ToLower(body.Password) {
		return echo.ErrUnauthorized
	}

	token, err := session.NewToken(session.JWTClaims{
		UserID:   user.ID,
		UserName: user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 168).Unix(),
		},
	})

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token":   token,
		"user_id": user.ID,
	})
}
