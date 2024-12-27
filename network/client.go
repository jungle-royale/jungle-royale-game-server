package network

import (
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ClientId string

type Client struct {
	ID   ClientId
	conn *websocket.Conn
}

func NewClient(conn *websocket.Conn) *Client {
	id := uuid.New().String()
	return &Client{ClientId(id), conn}
}

func (client *Client) write(data []byte) error {
	if err := client.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		log.Printf("Failed to send GameInit message to client %s: %v", client.ID, err)
		return err
	}
	return nil
}

func (client *Client) close() {
	client.conn.Close()
}
