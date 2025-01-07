package network

type StartMessageRequest struct {
	GameID    string   `json:"roomId"`
	ClientIds []string `json:"clientIds"`
}
