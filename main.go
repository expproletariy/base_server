package main

import (
	"fmt"
	"github.com/expproletariy/base_server/chat"
	"github.com/expproletariy/base_server/handlers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	id, err := chat.Manager().CreateRoom("test")
	if err != nil {
		panic(err)
	}
	fmt.Printf("New room %s\n", id)
	e := echo.New()
	e.Use(middleware.Logger())
	api := e.Group("/api")
	api.GET("/rooms", handlers.GetRooms)

	wsGroup := e.Group("/ws")
	wsGroup.GET("/room", handlers.ConnectToRoom)

	e.Logger.Fatal(e.Start(":3000"))
}
