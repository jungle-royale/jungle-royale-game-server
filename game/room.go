package game

type Room interface {
	OnClient(client *Client)
	OnMessage(data []byte, id string)
}
