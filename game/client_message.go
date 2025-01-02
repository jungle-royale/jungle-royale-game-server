package game

type ClientMessage struct {
	MessageType int
	RoomId      GameId
	ClientId    ClientId
	Data        []byte
}

func NewClientMessage(
	messageType int,
	gameId GameId,
	clientId ClientId,
	data []byte,
) *ClientMessage {
	return &ClientMessage{messageType, gameId, clientId, data}
}
