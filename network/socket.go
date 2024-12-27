package network

type Socket interface {
	Send(data []byte, clientId string) error
	Broadcast(data []byte)
}
