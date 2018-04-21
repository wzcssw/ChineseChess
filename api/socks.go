package api

import (
	"ChineseChess/chess"
	"fmt"
	"sync"

	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
)

type ConnectionPool struct {
	sync.RWMutex
	connections map[*websocket.Conn]struct{}
}

// WebSocks连接
var connectionPool ConnectionPool

func init() {
	connectionPool = ConnectionPool{}
	connectionPool.connections = make(map[*websocket.Conn]struct{})
}

func WebsocksProcess(c echo.Context) error {
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
			var msg string
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
				func(connection *websocket.Conn) {
					connectionPool.Lock()
					delete(connectionPool.connections, connection)
					connectionPool.Unlock()
				}(ws)
				break
			} else {
				//////////////
				fmt.Println("得到客户端的消息:", msg)
				result := chess.MsgProcess(msg)
				//////////////

				err = sendMessageToAllPool(result)
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
