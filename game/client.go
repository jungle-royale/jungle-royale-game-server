package game

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ClientId int

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
	return &Client{
		mu:             sync.Mutex{},
		GameID:         gameId,
		serverClientId: serverClientId,
		conn:           conn,
		connMu:         sync.Mutex{},
	}
}

func (client *Client) write(data []byte) error {
	// log.Printf("write start %s", client.serverClientId)
	client.connMu.Lock()
	defer client.connMu.Unlock()
	// log.Println(len(data))
	err := client.conn.WriteMessage(websocket.BinaryMessage, data)
	// log.Printf("write end %s", client.serverClientId)
	return err
}

func (client *Client) close() {
	client.connMu.Lock()
	defer client.connMu.Unlock()
	client.conn.Close()
	client.conn = nil
}
