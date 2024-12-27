package network

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ClientId string

type Client struct {
	ID     ClientId
	RoomID RoomId
	conn   *websocket.Conn
}

func NewClient(roomId RoomId, conn *websocket.Conn) *Client {
	id := uuid.New().String()
	return &Client{
		ClientId(id),
		roomId,
		conn,
	}
}

func (client *Client) write(data []byte) error {
	return client.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (client *Client) close() {
	client.conn.Close()
}
