package network

import "jungle-royale/statistic"

type EndMessageRequest struct {
	GameID  string                   `json:"roomId"`
	GameLog []statistic.ClientRecord `json:"users"`
}
