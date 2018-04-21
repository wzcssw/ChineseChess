package main

import (
	"ChineseChess/api"
	"ChineseChess/chess"
	"ChineseChess/lib"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(Process)

	e.Static("/", "./public")
	e.GET("/ws", api.WebsocksProcess)
	chess.MainChessPool = make(map[string]*chess.Chess)
	// MainChess = chess.InitChess()
	e.Logger.Fatal(e.Start(":1323"))
}

// 临时的回话处理: cookies 顺序匹配
func Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// set-cookie
		readCookie, err := c.Cookie("chinese_chess_user_id")
		ramdomTokenGenerator := lib.RamdomTokenGenerator()
		if err != nil {
			cookie := new(http.Cookie)
			cookie.Name = "chinese_chess_user_id"
			cookie.Value = ramdomTokenGenerator
			cookie.Expires = time.Now().Add(24 * time.Hour)
			c.SetCookie(cookie)
		} else {
			ramdomTokenGenerator = readCookie.Value
		}
		fmt.Println("[chinese_chess_user_id]=", ramdomTokenGenerator, " is old key:", err == nil)
		return next(c)
	}
}
