package game

type GameClientMessage struct {
	MessageType int
	ID          string
	Data        []byte // TODO: generic으로 처리
}

func NewGameClientMessage(messageType int, id string, data []byte) *GameClientMessage {
	return &GameClientMessage{messageType, id, data}
}
