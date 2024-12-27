package network

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	MaxClientCount = 100
	Port           = "8000"
)

type Socket struct {
	clients              map[ClientId]*Client
	clientsMu            sync.Mutex
	clientChannel        chan *Client
	clientMessageChannel chan *ClientMessage
	onClientCaller       func(ClientId)
	onMessageCaller      func(ClientId, []byte)
}

func NewSocket() *Socket {
	socket := Socket{
		make(map[ClientId]*Client),
		sync.Mutex{},
		make(chan *Client, MaxClientCount),
		make(chan *ClientMessage, MaxClientCount),
		nil,
		nil,
	}

	go func() {
		for client := range socket.clientChannel {
			go socket.setClient(client)
		}
	}()

	// message event handler
	go func() {
		for message := range socket.clientMessageChannel {
			go socket.handleClientMessage(message)
		}
	}()

	return &socket
}

func (socket *Socket) Listen() {

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
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
		defer conn.Close() // close시 client 삭제 처리?

		newClient := NewClient(conn)
		socket.clientChannel <- newClient

		log.Printf("Client %s connected", newClient.ID)

		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Client %s disconnected: %v", newClient.ID, err)
				break
			}

			socket.clientMessageChannel <- NewClientMessage(messageType, newClient.ID, data)
		}
	})

	log.Printf("Server listened in port " + Port)

	if err := http.ListenAndServe(":"+Port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func (socket *Socket) setClient(client *Client) {
	socket.clientsMu.Lock()
	socket.clients[client.ID] = client
	socket.clientsMu.Unlock()

	if socket.onClientCaller != nil {
		socket.onClientCaller(ClientId(client.ID))
	}
}

func (socket *Socket) handleClientMessage(clientMessage *ClientMessage) {
	if clientMessage.MessageType == websocket.BinaryMessage && socket.onMessageCaller != nil {
		socket.onMessageCaller(clientMessage.ID, clientMessage.Data)
	}
}

func (socket *Socket) OnClient(handler func(ClientId)) {
	socket.onClientCaller = handler
}

func (socket *Socket) OnMessage(handler func(ClientId, []byte)) {
	socket.onMessageCaller = handler
}

func (socket *Socket) Send(data []byte, id ClientId) {
	socket.clients[id].write(data)
}

func (socket *Socket) Broadcast(data []byte) {
	socket.clientsMu.Lock()
	for id, client := range socket.clients {
		if err := client.write(data); err != nil {
			log.Printf("Failed to send message to client %s: %v", id, err)
			client.close()
			delete(socket.clients, id)
		}
	}
	socket.clientsMu.Unlock()
}
