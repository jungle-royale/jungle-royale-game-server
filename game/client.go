package game

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ClientId string

type Client struct {
	mu     sync.Mutex
	ID     ClientId
	GameID GameId
	conn   *websocket.Conn
	connMu sync.Mutex
}

func NewClient(gameId GameId, conn *websocket.Conn) *Client {
	id := uuid.New().String()
	return &Client{
		sync.Mutex{},
		ClientId(id),
		gameId,
		conn,
		sync.Mutex{},
	}
}

// interface로 분리해서, Game에서는 interface에만 의존하도록 분리
// Game manager에서 Implementation을 넣어주도록 함
func (client *Client) write(data []byte) error {
	client.connMu.Lock()
	defer client.connMu.Unlock()
	return client.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (client *Client) close() {
	client.conn.Close()
}
