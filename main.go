package main

import (
	"chineseXiangqi/chess"
	"chineseXiangqi/tools"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/net/websocket"
)

// MainMap 棋盘
var MainMap chess.Map

// MsgProcess 通过得到的websockets消息处理Block
func MsgProcess(msg string) {
	fmt.Println("得到客户端的消息:" + msg)
	msgObj := tools.Msg{}
	json.Unmarshal([]byte(msg), &msgObj)
	switch msgObj.Operate {
	case "INIT_GAME":
		MainMap = chess.InitInstance()
		fmt.Println("初始化游戏")
	case "MOVE":
		MainMap.Move(msgObj)
		// MainMap.Test()
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
				MsgProcess(msg)
			}
			// Write
			bytes, _ := json.Marshal(MainMap.Print())
			err = sendMessageToAllPool(string(bytes))
			if err != nil {
				c.Logger().Error(err)
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
	e.Static("/", "./public")
	e.GET("/ws", hello)
	MainMap = chess.InitInstance()
	e.Logger.Fatal(e.Start(":1323"))
}
