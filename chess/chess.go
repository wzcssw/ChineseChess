package chess

import (
	"ChineseChess/lib"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// 局
type Chess struct {
	ID      string
	Map     Map
	term    int // 0: red   1: blue  2: over
	RUserID string
	BUserID string
}

// Point 坐标
type Point [2]int

// Map 棋盘
type Map map[string]Point
type reversedMap map[string]string

func InitChess() *Chess {
	chessObj := &Chess{}
	chessObj.Map = InitMap()
	chessObj.ID = lib.RamdomTokenGenerator()
	return chessObj
}

func (cs *Chess) Print() string {
	if cs == nil {
		return ""
	}
	result := make(map[string]interface{})
	result["map"] = cs.Map.Print()
	result["term"] = cs.term
	result["r_user_id"] = cs.RUserID
	result["b_user_id"] = cs.BUserID

	bytes, _ := json.Marshal(result)
	return string(bytes)
}

func (cs *Chess) TermString() string {
	if cs.term == 0 {
		return "R"
	} else {
		return "B"
	}
}

// InitMap 生成对象
func InitMap() Map {
	chess := make(map[string]Point)
	// 将军
	chess["BJiangJun"] = Point{4, 9}
	chess["RJiangJun"] = Point{4, 0}
	// 士
	chess["RShi1"] = Point{3, 0}
	chess["RShi2"] = Point{5, 0}
	chess["BShi1"] = Point{3, 9}
	chess["BShi2"] = Point{5, 9}
	// 象
	chess["RXiang1"] = Point{2, 0}
	chess["RXiang2"] = Point{6, 0}
	chess["BXiang1"] = Point{2, 9}
	chess["BXiang2"] = Point{6, 9}
	// 马
	chess["RMa1"] = Point{1, 0}
	chess["RMa2"] = Point{7, 0}
	chess["BMa1"] = Point{1, 9}
	chess["BMa2"] = Point{7, 9}
	// 车
	chess["RJu1"] = Point{0, 0}
	chess["RJu2"] = Point{8, 0}
	chess["BJu1"] = Point{0, 9}
	chess["BJu2"] = Point{8, 9}
	// 炮
	chess["RPao1"] = Point{1, 2}
	chess["RPao2"] = Point{7, 2}
	chess["BPao1"] = Point{1, 7}
	chess["BPao2"] = Point{7, 7}
	// 红卒
	chess["RZu1"] = Point{0, 3}
	chess["RZu2"] = Point{2, 3}
	chess["RZu3"] = Point{4, 3}
	chess["RZu4"] = Point{6, 3}
	chess["RZu5"] = Point{8, 3}
	// 蓝卒
	chess["BZu1"] = Point{0, 6}
	chess["BZu2"] = Point{2, 6}
	chess["BZu3"] = Point{4, 6}
	chess["BZu4"] = Point{6, 6}
	chess["BZu5"] = Point{8, 6}
	return chess
}

// Move 移动棋子
func (mp *Map) Move(msg lib.Msg) (lib.Msg, error) {
	originChess := ""
	for k, v := range *mp { // 找到对应棋子
		str := "[" + strconv.Itoa(v[0]) + "," + strconv.Itoa(v[1]) + "]"
		if str == msg.Origin {
			originChess = k
			break
		}
	}

	if originChess == "" { // 如果选中目标为空
		return msg, errors.New("选中目标为空")
	}

	targetPoint := ConvertStringToPoint(msg.Target)
	originPoint := (*mp)[originChess]
	if mp.valid(originChess, originPoint, targetPoint) {
		mp.deleteChess(targetPoint) // 如果targetPoint有棋子则删除
		(*mp)[originChess] = targetPoint
	} else {
		return msg, errors.New("走法违规")
	}
	return msg, nil
}

// Move 移动棋子
func (cs *Chess) Move(msg lib.Msg) lib.Msg {
	// cs.term
	Msg, err := cs.Map.Move(msg)
	if err == nil {
		if cs.RUserID == msg.UserID {
			cs.term = 1
		} else {
			cs.term = 0
		}
	}

	return Msg
}

// deleteChess 如果targetPoint有棋子则删除
func (mp *Map) deleteChess(targetPoint Point) bool {
	result := false
	for k, v := range *mp {
		if (v[0] == targetPoint[0]) && (v[1] == targetPoint[1]) {
			delete(*mp, k)
			result = true
			break
		}
	}
	return result
}

// valid 验证是否可走
func (mp *Map) valid(originChess string, originPoint Point, targetPoint Point) bool {
	result := false
	if mp.checkWChess(string(originChess[0]), targetPoint) { // 判断targetPoint是否有友军
		return false
	}
	if strings.Contains(originChess, "JiangJun") {
		result = validJiangJun(mp, originChess, originPoint, targetPoint)
	} else if strings.Contains(originChess, "Shi") {
		result = validShi(originChess, originPoint, targetPoint)
	} else if strings.Contains(originChess, "Xiang") {
		result = validXiang(mp, originChess, originPoint, targetPoint)
	} else if strings.Contains(originChess, "Ma") {
		result = validMa(mp, originChess, originPoint, targetPoint)
	} else if strings.Contains(originChess, "Ju") {
		result = validJu(mp, originChess, originPoint, targetPoint)
	} else if strings.Contains(originChess, "Pao") {
		result = validPao(mp, originChess, originPoint, targetPoint)
	} else if strings.Contains(originChess, "Zu") {
		result = validZu(originChess, originPoint, targetPoint)
	} else {
	}
	return result
}

// 将军
func validJiangJun(mp *Map, originChess string, originPoint Point, targetPoint Point) bool {
	result := false
	if strings.Contains(mp.getChessName(targetPoint), "JiangJun") { // 如果目标是将军
		countChess := mp.countLineChess(originPoint, targetPoint) == 0 // 线上的棋子数量为0
		sameLine := originPoint[0] == targetPoint[0]                   // 在一条线上
		if countChess && sameLine {
			return true
		}
	}
	if string(originChess[0]) == "B" { // 蓝棋
		c1 := (targetPoint[0] > 2 && targetPoint[0] < 6)
		c2 := targetPoint[1] > 6
		c3 := (abs(targetPoint[0]-originPoint[0]) == 1) && (targetPoint[1] == originPoint[1]) // x走一步且y不动
		c4 := (abs(targetPoint[1]-originPoint[1]) == 1) && (targetPoint[0] == originPoint[0]) // y走一步且x不动
		if (c1 && c2) && (c3 || c4) {
			result = true
		}
	} else { // 红棋
		c1 := (targetPoint[0] > 2 && targetPoint[0] < 6)
		c2 := targetPoint[1] < 3
		c3 := (abs(targetPoint[0]-originPoint[0]) == 1) && (targetPoint[1] == originPoint[1]) // x走一步且y不动
		c4 := (abs(targetPoint[1]-originPoint[1]) == 1) && (targetPoint[0] == originPoint[0]) // y走一步且x不动
		if (c1 && c2) && (c3 || c4) {
			result = true
		}
	}
	return result
}

// 士
func validShi(originChess string, originPoint Point, targetPoint Point) bool {
	result := false
	if string(originChess[0]) == "B" { // 蓝棋
		c1 := (targetPoint[0] > 2 && targetPoint[0] < 6)
		c2 := targetPoint[1] > 6
		c3 := ((abs(targetPoint[0]-originPoint[0]) == 1) && (abs(targetPoint[1]-originPoint[1]) == 1)) // x,y都只走一步
		if c1 && c2 && c3 {
			result = true
		}
	} else { // 红棋
		c1 := (targetPoint[0] > 2 && targetPoint[0] < 6)
		c2 := targetPoint[1] < 3
		c3 := ((abs(targetPoint[0]-originPoint[0]) == 1) && (abs(targetPoint[1]-originPoint[1]) == 1)) // x,y都只走一步
		if c1 && c2 && c3 {
			result = true
		}
	}
	return result
}

// 象
func validXiang(mp *Map, originChess string, originPoint Point, targetPoint Point) bool {
	result := false
	if string(originChess[0]) == "B" { // 蓝棋
		c1 := targetPoint[1] > 4
		c2 := ((abs(targetPoint[0]-originPoint[0]) == 2) && (abs(targetPoint[1]-originPoint[1]) == 2)) // x,y都走2步
		if c1 && c2 {
			result = true
		}
	} else { // 红棋
		c1 := targetPoint[1] < 5
		c2 := ((abs(targetPoint[0]-originPoint[0]) == 2) && (abs(targetPoint[1]-originPoint[1]) == 2)) // x,y都走2步
		if c1 && c2 {
			result = true
		}
	}
	// 别腿
	xDiff := targetPoint[0] - originPoint[0]
	yDiff := targetPoint[1] - originPoint[1]
	centerPoint := Point{originPoint[0] + xDiff/2, originPoint[1] + yDiff/2}
	if result {
		result = !mp.checkChess(centerPoint)
	}
	return result
}

// 马
func validMa(mp *Map, originChess string, originPoint Point, targetPoint Point) bool {
	result := false
	c1 := ((abs(targetPoint[0]-originPoint[0]) == 1) && (abs(targetPoint[1]-originPoint[1]) == 2)) // x走1步,y走2步
	c2 := ((abs(targetPoint[0]-originPoint[0]) == 2) && (abs(targetPoint[1]-originPoint[1]) == 1)) // x走2步,y走1步
	if (c1 && !c2) || (!c1 && c2) {
		result = true
	}
	// 别腿
	xDiff := targetPoint[0] - originPoint[0]
	yDiff := targetPoint[1] - originPoint[1]
	if result {
		if abs(xDiff) == 2 {
			// 如果是横着走
			result = !mp.checkChess(Point{targetPoint[0] - xDiff/2, originPoint[1]})
		} else if abs(yDiff) == 2 {
			// 如果是竖着走
			result = !mp.checkChess(Point{originPoint[0], targetPoint[1] - yDiff/2})
		}
	}

	return result
}

// 车
func validJu(mp *Map, originChess string, originPoint Point, targetPoint Point) bool {
	result := false
	c1 := (originPoint[0] == targetPoint[0])                 // x轴相等
	c2 := (originPoint[1] == targetPoint[1])                 // y轴相等
	c3 := (mp.countLineChess(originPoint, targetPoint) == 0) // 路径上是否有其他棋子
	if ((c1 && !c2) || (!c1 && c2)) && c3 {
		result = true
	}
	return result
}

// 炮
func validPao(mp *Map, originChess string, originPoint Point, targetPoint Point) bool {
	result := false
	countChess := mp.countLineChess(originPoint, targetPoint)
	c1 := (originPoint[0] == targetPoint[0]) // x轴相等
	c2 := (originPoint[1] == targetPoint[1]) // y轴相等

	armyColor := "B"
	if string(originChess[0]) == "B" {
		armyColor = "R"
	}
	c4 := mp.checkWChess(armyColor, targetPoint) // 判断targetPoint是否有敌军

	if c4 { // 吃棋
		c3 := (countChess == 1) // 路径上只有一个棋子
		result = ((c1 && !c2) || (!c1 && c2)) && c3
	} else { // 移动
		c5 := !mp.checkChess(targetPoint) // 判断targetPoint是棋子
		result = ((c1 && !c2) || (!c1 && c2)) && (countChess == 0) && c5
	}
	return result
}

// 卒子
func validZu(originChess string, originPoint Point, targetPoint Point) bool {
	result := false
	if string(originChess[0]) == "B" { // 蓝棋
		if originPoint[1] > 4 { // 在河内
			if (targetPoint[0] == originPoint[0]) && ((targetPoint[1] + 1) == originPoint[1]) {
				result = true
			}
		} else { //在河外
			c1 := targetPoint[1] <= originPoint[1]                                                      // 不可以倒着走
			c2 := ((targetPoint[1] + 1) == originPoint[1]) && (abs(targetPoint[0]-originPoint[0]) == 0) // y轴走一步且没有左右走
			c3 := ((abs(targetPoint[0]-originPoint[0]) == 1) && (targetPoint[1] == originPoint[1]))     //   左右走一步且没有前后走
			if c1 && (c2 || c3) {
				result = true
			}
		}
	} else { // 红棋
		if originPoint[1] < 5 { // 在河内
			if (targetPoint[0] == originPoint[0]) && ((targetPoint[1] - 1) == originPoint[1]) {
				result = true
			}
		} else { //在河外
			c1 := targetPoint[1] >= originPoint[1]                                                      // 不可以倒着走
			c2 := ((targetPoint[1] - 1) == originPoint[1]) && (abs(targetPoint[0]-originPoint[0]) == 0) // y轴走一步且没有左右走
			c3 := ((abs(targetPoint[0]-originPoint[0]) == 1) && (targetPoint[1] == originPoint[1]))     //   左右走一步且没有前后走
			if c1 && (c2 || c3) {
				result = true
			}
		}
	}
	return result
}

// Print 打印输出 x:0~8  y:0~9
func (mp *Map) Print() map[string]string {
	reversedMap := mp.reverseMap()
	result := make(map[string]string)
	for x := 0; x <= 8; x++ {
		for y := 0; y <= 9; y++ {
			pointStr := "[" + strconv.Itoa(x) + "," + strconv.Itoa(y) + "]"
			chessName := reversedMap[pointStr]
			result[pointStr] = chessName
		}
	}
	return result
}

// reverseMap 将Map转化为reversedMap
func (mp *Map) reverseMap() reversedMap {
	reversedMap := make(reversedMap)
	for k, v := range *mp {
		pointStr := "[" + strconv.Itoa(v[0]) + "," + strconv.Itoa(v[1]) + "]"
		reversedMap[pointStr] = k
	}
	return reversedMap
}

// ConvertStringToPoint string => Point
func ConvertStringToPoint(params string) Point {
	point := Point{}
	params = strings.Replace(params, "[", "", -1)
	params = strings.Replace(params, "]", "", -1)
	strArray := strings.Split(params, ",")
	for i, s := range strArray {
		point[i], _ = strconv.Atoi(s)
	}
	return point
}

// abs 绝对值
func abs(a int) (ret int) {
	ret = (a ^ a>>31) - a>>31
	return
}

// Test test
func (mp *Map) Test() {
	for k, v := range *mp {
		fmt.Println(k, v)
	}
}

// checkChess 检查该位置是否有(_颜色_)的棋子
func (mp *Map) checkWChess(color string, targetPoint Point) bool {
	result := false
	for k, v := range *mp {
		if (v[0] == targetPoint[0]) && (v[1] == targetPoint[1]) {
			if string(k[0]) == color {
				result = true
				break
			}
		}
	}
	return result
}

// checkChess 检查该位置是否有棋子
func (mp *Map) checkChess(targetPoint Point) bool {
	result := false
	for _, v := range *mp {
		if (v[0] == targetPoint[0]) && (v[1] == targetPoint[1]) {
			result = true
			break
		}
	}
	return result
}

// checkChess 检查该位置是否有棋子
func (mp *Map) getChessName(targetPoint Point) string {
	result := ""
	for k, v := range *mp {
		if (v[0] == targetPoint[0]) && (v[1] == targetPoint[1]) {
			result = k
			break
		}
	}
	return result
}

// checkChess 两点间的棋子数量
func (mp *Map) countLineChess(originPoint Point, targetPoint Point) int {
	result := 0
	if originPoint[0] == targetPoint[0] { // y轴
		yDiff := targetPoint[1] - originPoint[1]
		if yDiff > 0 {
			for i := (originPoint[1] + 1); i < targetPoint[1]; i++ {
				if mp.checkChess(Point{originPoint[0], i}) {
					result++
				}
			}
		} else {
			for j := (targetPoint[1] + 1); j < originPoint[1]; j++ {
				if mp.checkChess(Point{originPoint[0], j}) {
					result++
				}
			}
		}
	} else if originPoint[1] == targetPoint[1] { // x轴
		xDiff := targetPoint[0] - originPoint[0]
		if xDiff > 0 {
			for i := (originPoint[0] + 1); i < targetPoint[0]; i++ {
				if mp.checkChess(Point{i, originPoint[1]}) {
					result++
				}
			}
		} else {
			for j := (targetPoint[0] + 1); j < originPoint[0]; j++ {
				if mp.checkChess(Point{j, originPoint[1]}) {
					result++
				}
			}
		}
	}
	return result
}
