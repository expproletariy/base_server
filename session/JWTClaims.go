package session

import "github.com/dgrijalva/jwt-go"

type JWTClaims struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.StandardClaims
}
