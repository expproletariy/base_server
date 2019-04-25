package types

import "time"

type SimpleMessage struct {
	ID        string    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Text      string    `json:"text" db:"text"`
	UserID    string    `json:"user_id" db:"users_id"`
	UserName  string    `json:"user_name" db:"users_name"`
}
