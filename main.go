package main

import (
	"ChineseChess/api"
	"ChineseChess/lib"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(lib.SessionProcess)

	e.Static("/", "./public")
	e.GET("/ws", api.WebsocksProcess)

	e.Logger.Fatal(e.Start(":1323"))
}
