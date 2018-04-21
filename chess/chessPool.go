package chess

import (
	"ChineseChess/lib"
	"encoding/json"
	"fmt"
)

type ChessPool map[string]*Chess

// MainChess 棋盘
var MainChessPool ChessPool

func (cp *ChessPool) GetUsersChess(userID string) *Chess {
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

func (cp *ChessPool) ReadyChess(userID string) *Chess {
	//  查找userID已经参与的Chess
	chessID := ""
	for _, chess := range *cp {
		if chess.BUserID == "" || chess.RUserID == "" {
			chessID = chess.ID
		}
	}
	return (*cp)[chessID]
}

func (cp *ChessPool) InitNewChess(userID string) *Chess {
	newChess := InitChess()
	newChess.RUserID = userID
	(*cp)[newChess.ID] = newChess
	return newChess
}

// 消息处理
func MsgProcess(msg string) string {
	msgObj := lib.Msg{}
	json.Unmarshal([]byte(msg), &msgObj)

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

	return MainChessPool.GetUsersChess(msgObj.UserID).Print()
}
