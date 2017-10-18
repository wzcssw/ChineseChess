package tools

// Msg 块
type Msg struct {
	PlayerName string `json:"player_name"`
	Operate    string `json:"operate"`
	Origin     string `json:"origin"`
	Target     string `json:"target"`
}
