package network

type GameServerNotificationRequest struct {
	GameID      string `json:"roomId"`
	MinPlayers  int    `json:"minPlayers"`
	MaxPlayTime int    `json:"maxPlayTime"` // seconds
}
