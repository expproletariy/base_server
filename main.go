package main

import (
	"github.com/expproletariy/base_server/handlers"
	"github.com/expproletariy/base_server/models"
	"github.com/expproletariy/base_server/session"
	"github.com/expproletariy/base_server/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
)

var PORT string
var DbUser string
var DbPass string
var DbHost string
var DbPort string
var DbName string
var DSN string

func init() {
	DbUser = os.Getenv("DB_USER")
	DbPass = os.Getenv("DB_PASS")
	DbHost = os.Getenv("DB_HOST")
	DbPort = os.Getenv("DB_PORT")
	DbName = os.Getenv("DB_NAME")
	PORT = os.Getenv("PORT")

	//Build db connection string in DSN
	DSN = DbUser + ":" + DbPass + "@tcp(" + DbHost + ":" + DbPort + ")/" + DbName + "?parseTime=true"
	conn, err := sqlx.Connect("mysql", DSN)
	if err != nil {
		panic(err)
	}
	models.SetContext(conn)
	err = models.Prepare()
	if err != nil {
		panic(err)
	}
}

func main() {
	if ctx, ok := models.GetContext(); ok {
		defer ctx.Close()
	} else {
		panic(types.NewError("Empty db context, before use need to SetContext"))
	}
	e := echo.New()
	e.Use(middleware.Logger())
	api := e.Group("/api")
	api.POST("/login", handlers.Login)
	api.POST("/sign_in", handlers.SignIn)
	v1 := api.Group("/v1")
	v1.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     new(session.JWTClaims),
		SigningKey: session.SecretKey(),
	}))
	v1.GET("/rooms", handlers.GetRooms)
	v1.POST("/room/create", handlers.CreateRoom)
	v1.POST("/user/block", handlers.BlockUser)
	v1.GET("/user/messages", handlers.GetMessageHistory)

	wsGroup := e.Group("/ws")
	v1.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     new(session.JWTClaims),
		SigningKey: session.SecretKey(),
	}))
	wsGroup.GET("/room", handlers.ConnectToRoom)

	e.Logger.Fatal(e.Start(PORT))
}
