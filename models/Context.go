package models

import (
	"github.com/jmoiron/sqlx"
)

var context *sqlx.DB

//SetContext bind context
func SetContext(db *sqlx.DB) {
	context = db
}

//GetContext of db connection
func GetContext() (*sqlx.DB, bool) {
	if context == nil {
		return nil, false
	}
	return context, true
}
