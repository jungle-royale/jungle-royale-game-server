package network

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ClientId string

type Client struct {
	ID     ClientId
	RoomID RoomId
	conn   *websocket.Conn
	connMu sync.Mutex
}

func NewClient(roomId RoomId, conn *websocket.Conn) *Client {
	id := uuid.New().String()
	return &Client{
		ClientId(id),
		roomId,
		conn,
		sync.Mutex{},
	}
}

func (client *Client) write(data []byte) error {
	client.connMu.Lock()
	defer client.connMu.Unlock()
	return client.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (client *Client) close() {
	client.conn.Close()
}
