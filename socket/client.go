package socket

import "github.com/gorilla/websocket"

type Client struct {
	conn *websocket.Conn
	id   string
}

func NewClient(conn *websocket.Conn, id string) *Client {
	return &Client{conn, id}
}
