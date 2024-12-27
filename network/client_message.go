package network

type ClientMessage struct {
	MessageType int
	RoomId      RoomId
	ClientId    ClientId
	Data        []byte
}

func NewClientMessage(
	messageType int,
	roomId RoomId,
	clientId ClientId,
	data []byte,
) *ClientMessage {
	return &ClientMessage{messageType, roomId, clientId, data}
}
