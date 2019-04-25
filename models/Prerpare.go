package models

import (
	"github.com/expproletariy/base_server/types"
	"github.com/jmoiron/sqlx"
)

var (
	StmtGetRooms                    *sqlx.Stmt
	StmtGetRoom                     *sqlx.Stmt
	StmtGetUser                     *sqlx.Stmt
	StmtGetUserByName               *sqlx.Stmt
	StmtGetUserRoom                 *sqlx.Stmt
	StmtCheckUserRoom               *sqlx.Stmt
	StmtCheckRoomName               *sqlx.Stmt
	StmtSaveMessage                 *sqlx.NamedStmt
	StmtNewUserRoomByDefault        *sqlx.NamedStmt
	StmtNewUser                     *sqlx.NamedStmt
	StmtUserBlocker                 *sqlx.Stmt
	StmtGetMessageHistory           *sqlx.Stmt
	StmtGetMessageHistoryForBlocked *sqlx.Stmt
)

func Prepare() error {
	if ctx, ok := GetContext(); ok {
		var err error
		StmtGetRooms, err = ctx.Preparex("SELECT * FROM rooms")
		if err != nil {
			return err
		}
		StmtGetRoom, err = ctx.Preparex("SELECT * FROM rooms WHERE rooms.id=? LIMIT 1")
		if err != nil {
			return err
		}
		StmtGetUser, err = ctx.Preparex("SELECT * FROM users WHERE users.id=? LIMIT 1")
		if err != nil {
			return err
		}
		StmtGetUserByName, err = ctx.Preparex("SELECT * FROM users WHERE users.name=? LIMIT 1")
		if err != nil {
			return err
		}
		StmtGetUserRoom, err = ctx.Preparex("SELECT " +
			"users.id, users.name, user_room.blocked, user_room.creator " +
			"FROM users INNER JOIN user_room ON user_room.users_id=users.id " +
			"WHERE user_room.users_id=? AND user_room.rooms_id=? LIMIT 1",
		)
		if err != nil {
			return err
		}
		StmtSaveMessage, err = ctx.PrepareNamed("INSERT INTO messages " +
			"(id, text, users_id, rooms_id, created_at) " +
			"VALUES (:id, :text, :users_id, :rooms_id, :created_at)",
		)
		if err != nil {
			return err
		}
		StmtNewUserRoomByDefault, err = ctx.PrepareNamed("INSERT INTO user_room " +
			"(users_id, rooms_id) " +
			"VALUES (:users_id, :rooms_id)",
		)
		if err != nil {
			return err
		}
		StmtNewUser, err = ctx.PrepareNamed("INSERT INTO users " +
			"(id, name, password) " +
			"VALUES (:id, :name, :password)",
		)
		if err != nil {
			return err
		}
		StmtCheckRoomName, err = ctx.Preparex("SELECT name FROM rooms WHERE name=? LIMIT 1")
		if err != nil {
			return err
		}
		//args: blocked, users_id, rooms_id
		StmtUserBlocker, err = ctx.Preparex("UPDATE user_room SET blocked=?, block_time=? WHERE users_id=? AND rooms_id=?")
		if err != nil {
			return err
		}
		//args: users_id, rooms_id, page
		StmtGetMessageHistory, err = ctx.Preparex(`SELECT ` +
			`messages.id AS "id",` +
			`messages.text AS "text",` +
			`user_room.users_id AS "users_id",` +
			`users.name AS "users_name",` +
			`messages.created_at AS "created_at" ` +
			`FROM messages ` +
			`INNER JOIN user_room ON user_room.users_id=messages.users_id AND user_room.rooms_id=messages.rooms_id ` +
			`INNER JOIN users ON user_room.users_id=users.id ` +
			`WHERE messages.rooms_id=? ` +
			`ORDER BY messages.created_at ` +
			`LIMIT ?, 20`,
		)
		if err != nil {
			return err
		}
		//args: users_id, rooms_id, page
		StmtGetMessageHistoryForBlocked, err = ctx.Preparex(`SELECT ` +
			`messages.id AS "id",` +
			`messages.text AS "text",` +
			`user_room.users_id AS "users_id",` +
			`users.name AS "users_name",` +
			`messages.created_at AS "created_at" ` +
			`FROM messages ` +
			`INNER JOIN user_room ON user_room.users_id=messages.users_id AND user_room.rooms_id=messages.rooms_id ` +
			`INNER JOIN users ON user_room.users_id=users.id ` +
			`WHERE messages.users_id=? AND messages.rooms_id=? AND messages.created_at < user_room.block_time ` +
			`ORDER BY messages.created_at ` +
			`LIMIT ?, 20`,
		)
		if err != nil {
			return err
		}
		return nil
	}
	return types.NewError("Empty db context, before use need to SetContext")
}
