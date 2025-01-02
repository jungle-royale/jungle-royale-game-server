package network

type GameServerNotificationRequest struct {
	RoomID      string `json:"roomId"`
	MinPlayers  int    `json:"minPlayers"`
	MaxPlayTime int    `json:"maxPlayTime"` // seconds
}
