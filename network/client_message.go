package network

type ClientMessage struct {
	MessageType int
	ID          ClientId
	Data        []byte
}

func NewClientMessage(messageType int, id ClientId, data []byte) *ClientMessage {
	return &ClientMessage{messageType, id, data}
}
