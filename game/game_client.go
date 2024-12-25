package game

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type GameClient struct {
	ID   string
	Conn *websocket.Conn
}

func NewGameClient(conn *websocket.Conn) *GameClient {
	id := uuid.New().String()
	return &GameClient{id, conn}
}
