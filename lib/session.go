package lib

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

// 临时的回话处理: cookies 顺序匹配
func SessionProcess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// set-cookie
		readCookie, err := c.Cookie("chinese_chess_user_id")
		ramdomTokenGenerator := RamdomTokenGenerator()
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
