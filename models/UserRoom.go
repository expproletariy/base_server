package models

import "github.com/expproletariy/base_server/types"

type UserRoom struct {
	ID     uint64 `json:"id" db:"id"`
	RoomID string `json:"room_id" db:"rooms_id"`
	UserID string `json:"user_id" db:"users_id"`
	types.User
}
