package network

type PlayerLeaveMessageRequest struct {
	GameID   string `json:"roomId"`
	ClientID string `json:"clientId"`
}
