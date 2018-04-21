package lib

// 返回消息
type Msg struct {
	UserID     string `json:"user_id"`
	PlayerName string `json:"player_name"`
	Operate    string `json:"operate"`
	Origin     string `json:"origin"`
	Target     string `json:"target"`
}
