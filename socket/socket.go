package socket

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
}

type ClientMessage struct {
	MessageType int
	ID          string
	Data        []byte
}

var ClientChannel chan *Client
var MessageChannel chan *ClientMessage

func NewClient(conn *websocket.Conn) *Client {
	id := generateClientID()
	return &Client{id, conn}
}

func NewClientMessage(messageType int, ID string, data []byte) *ClientMessage {
	return &ClientMessage{messageType, ID, data}
}

func generateClientID() string {
	return uuid.New().String()
}

func InitSocket() {

	ClientChannel = make(chan *Client, 100)
	MessageChannel = make(chan *ClientMessage, 100)

	http.HandleFunc("/ws", handleWebSocket)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// http → websocket upgrade
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	newClient := NewClient(conn)

	ClientChannel <- newClient

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Client %s disconnected: %v", newClient.ID, err)
			break
		}

		MessageChannel <- NewClientMessage(messageType, newClient.ID, data)
	}
}
