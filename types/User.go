package types

import "time"

type User struct {
	ID        string     `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Blocked   bool       `json:"blocked" db:"blocked"`
	BlockTime *time.Time `json:"block_time" db:"block_time"`
	Creator   bool       `json:"creator" db:"creator"`
}
