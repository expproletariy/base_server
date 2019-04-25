package session

import (
	"github.com/dgrijalva/jwt-go"
	uuid2 "github.com/satori/go.uuid"
	"sync"
)

var secretKey []byte
var once sync.Once

func SecretKey() []byte {
	once.Do(func() {
		secretKey = []byte(uuid2.NewV4().String())
	})
	return secretKey
}

func NewToken(claims JWTClaims) (string, error) {
	return jwt.
		NewWithClaims(jwt.SigningMethodHS256, &claims).
		SignedString(SecretKey())
}

func GetClaimsInfo(auth string) (*JWTClaims, bool) {

	if len(auth) == 0 {
		return nil, false
	}

	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(auth, claims, func(token *jwt.Token) (i interface{}, e error) {
		return SecretKey(), nil
	})

	if err != nil {
		return nil, false
	}
	if !token.Valid {
		return nil, false
	}

	return token.Claims.(*JWTClaims), true
}
