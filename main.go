package main

import (
	"ChineseChess/chess"
	"ChineseChess/tools"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/net/websocket"
)

type ChessPool map[string]*chess.Chess

// MainChess 棋盘
var MainChessPool ChessPool

func (cp *ChessPool) GetUsersChess(userID string) *chess.Chess {
	//  查找userID已经参与的Chess
	chessID := ""
	for _, chess := range *cp {
		if chess.RUserID == userID || chess.BUserID == userID {
			chessID = chess.ID
		}
	}
	// 如果未找到匹配Chess则匹配一个正在等待方
	if chessID == "" {
		if readyChess := cp.ReadyChess(userID); readyChess == nil { // 如果没有等待者返回新chess
			return cp.InitNewChess(userID)
		} else {
			readyChess.BUserID = userID
			return readyChess
		}
	}
	return (*cp)[chessID]
}

func (cp *ChessPool) ReadyChess(userID string) *chess.Chess {
	//  查找userID已经参与的Chess
	chessID := ""
	for _, chess := range *cp {
		if chess.BUserID == "" || chess.RUserID == "" {
			chessID = chess.ID
		}
	}
	return (*cp)[chessID]
}

func (cp *ChessPool) InitNewChess(userID string) *chess.Chess {
	newChess := chess.InitChess()
	newChess.RUserID = userID
	(*cp)[newChess.ID] = newChess
	return newChess
}

func MsgProcess(msgObj tools.Msg) {
	switch msgObj.Operate {
	case "INIT_GAME":
		fmt.Println("初始化游戏")
	case "MOVE":
		chess := MainChessPool.GetUsersChess(msgObj.UserID)
		chess.Move(msgObj)
	case "RELOAD_GAME":
	default:
		fmt.Printf("Default Default")
	}
}

// 连接池
var connectionPool = struct {
	sync.RWMutex
	connections map[*websocket.Conn]struct{}
}{
	connections: make(map[*websocket.Conn]struct{}),
}

func hello(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			connectionPool.Lock()
			connectionPool.connections[ws] = struct{}{}

			defer func(connection *websocket.Conn) {
				connectionPool.Lock()
				delete(connectionPool.connections, connection)
				connectionPool.Unlock()
			}(ws)
			connectionPool.Unlock()
			// Read
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				fmt.Printf("===========>In Error: %s\n", err)
				c.Logger().Error(err)
				func(connection *websocket.Conn) {
					connectionPool.Lock()
					delete(connectionPool.connections, connection)
					connectionPool.Unlock()
				}(ws)
				break
			} else {
				fmt.Println("得到客户端的消息:", msg)
				msgObj := tools.Msg{}
				json.Unmarshal([]byte(msg), &msgObj)
				MsgProcess(msgObj)
				// Write
				fmt.Println("棉花糖", MainChessPool.GetUsersChess(msgObj.UserID))
				bytes, _ := json.Marshal(MainChessPool.GetUsersChess(msgObj.UserID).Print())
				err = sendMessageToAllPool(string(bytes))
				if err != nil {
					c.Logger().Error(err)
				}
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

// 广播
func sendMessageToAllPool(message string) error {
	connectionPool.RLock()
	defer connectionPool.RUnlock()
	for connection := range connectionPool.connections {
		if err := websocket.Message.Send(connection, message); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(Process)

	e.Static("/", "./public")
	e.GET("/ws", hello)
	MainChessPool = make(map[string]*chess.Chess)
	// MainChess = chess.InitChess()
	e.Logger.Fatal(e.Start(":1323"))
}

// 临时的回话处理: cookies 顺序匹配
func Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		/////// set-cookie
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
		///////
		return next(c)
	}
}

//  生成随机字符串
func RamdomTokenGenerator() string {
	rand.Seed(time.Now().UnixNano())
	x := rand.Intn(1e11)
	data := []byte(strconv.Itoa(x))
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
