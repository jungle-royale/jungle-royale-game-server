package game

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type ClientId int

type Client struct {
	mu             sync.Mutex
	isClose        bool
	ID             ClientId
	GameID         GameId
	serverClientId string
	userName       string
	conn           *websocket.Conn
	sendChan       chan []byte
	isObserver     bool
	closeOnce      *sync.Once
}

func NewClient(
	gameId GameId,
	serverClientId string,
	userName string,
	conn *websocket.Conn,
	isObserver bool,
	clientId int,
) *Client {
	newClient := &Client{
		mu:             sync.Mutex{},
		isClose:        false,
		ID:             ClientId(clientId),
		GameID:         gameId,
		serverClientId: serverClientId,
		userName:       userName,
		conn:           conn,
		sendChan:       make(chan []byte, 200),
		isObserver:     isObserver,
		closeOnce:      new(sync.Once),
	}
	go newClient.SendData()
	return newClient
}

func (client *Client) write(data []byte) {
	if client.isClose {
		return
	}
	select {
	case client.sendChan <- data:
	default:
		log.Printf("Send channel is full, dropping message for client %s", client.serverClientId)
	}
}

func (client *Client) SendData() {
	for data := range client.sendChan {
		// log.Println(len(data))
		if client.conn == nil {
			continue
		}
		err := client.conn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			log.Printf("err while sending data to client %s: %v", client.serverClientId, err)
		}
	}

}

func (client *Client) close() {

	client.closeOnce.Do(func() {
		client.isClose = true

		if client.conn != nil {
			client.conn.Close()
		}

		if client.sendChan != nil {
			close(client.sendChan)
		}

		client.conn = nil
	})
}
