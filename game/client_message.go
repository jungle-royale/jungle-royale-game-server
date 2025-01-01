package game

type ClientMessage struct {
	MessageType int
	RoomId      GameId
	ClientId    ClientId
	Data        []byte
}

func NewClientMessage(
	messageType int,
	roomId GameId,
	clientId ClientId,
	data []byte,
) *ClientMessage {
	return &ClientMessage{messageType, roomId, clientId, data}
}
