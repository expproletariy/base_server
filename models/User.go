package models

import "github.com/expproletariy/base_server/types"

type User struct {
	Password string `json:"password" db:"password"`
	types.User
}
