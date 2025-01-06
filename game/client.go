package game

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ClientId string

type Client struct {
	mu             sync.Mutex
	ID             ClientId
	GameID         GameId
	serverClientId string
	conn           *websocket.Conn
	connMu         sync.Mutex
}

func NewClient(
	gameId GameId,
	serverClientId string,
	conn *websocket.Conn,
) *Client {
	id := uuid.New().String()
	return &Client{
		mu:             sync.Mutex{},
		ID:             ClientId(id),
		GameID:         gameId,
		serverClientId: serverClientId,
		conn:           conn,
		connMu:         sync.Mutex{},
	}
}

func (client *Client) write(data []byte) error {
	client.connMu.Lock()
	defer client.connMu.Unlock()
	return client.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (client *Client) close() {
	client.conn.Close()
	client.conn = nil
}
