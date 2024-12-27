package network

type Room interface {
	OnClient(clientId string)
	OnMessage(data []byte, id string)
}
