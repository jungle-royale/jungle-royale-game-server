package game

import (
	"log"
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
	sendChan       chan []byte
	isObserver     bool
}

func NewClient(
	gameId GameId,
	serverClientId string,
	conn *websocket.Conn,
	isObserver bool,
) *Client {
	newClient := &Client{
		mu:             sync.Mutex{},
		GameID:         gameId,
		serverClientId: serverClientId,
		conn:           conn,
		sendChan:       make(chan []byte, 4),
		isObserver:     isObserver,
	}
	go newClient.SendData()
	return newClient
}

func (client *Client) write(data []byte) {
	select {
	case client.sendChan <- data:
	default:
		log.Printf("Send channel is full, dropping message for client %s", client.serverClientId)
	}
}

func (client *Client) SendData() {
	for data := range client.sendChan {
		// log.Println(len(data))
		err := client.conn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			log.Printf("err while sending data to client %s", client.serverClientId)
		}
	}
}

func (client *Client) close() {

	if client.conn != nil {
		client.conn.Close()
	}

	select {
	case <-client.sendChan:
	default:
		close(client.sendChan)
	}

	client.conn = nil
}
