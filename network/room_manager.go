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

type RoomId string // room create 할 때 받음

type RoomManager struct {
	rooms                map[RoomId]*Room
	roomsMu              sync.Mutex
	clients              map[ClientId]*Client
	clientsMu            sync.Mutex
	clientChannel        chan *Client
	clientMessageChannel chan *ClientMessage
}

// add game

func NewRoomManager() *RoomManager {
	socket := RoomManager{
		make(map[RoomId]*Room),
		sync.Mutex{},
		make(map[ClientId]*Client),
		sync.Mutex{},
		make(chan *Client, MaxClientCount),
		make(chan *ClientMessage, MaxClientCount),
	}
	return &socket
}

func (socket *RoomManager) Listen() {

	go func() {
		for client := range socket.clientChannel {
			go socket.setClient(client)
		}
	}()

	go func() {
		for message := range socket.clientMessageChannel {
			go socket.handleClientMessage(message)
		}
	}()

	http.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) {
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

		roomId := r.URL.Query().Get("roomId")
		if roomId == "" {
			http.Error(w, "Missing roomId query parameter", http.StatusBadRequest)
			return
		}

		log.Printf("Received connection request for Room ID: %s", roomId)

		newClient := NewClient(RoomId(roomId), conn)
		socket.clientChannel <- newClient

		log.Printf("Client %s connected", newClient.ID)

		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Client %s disconnected: %v", newClient.ID, err)
				break
			}

			socket.clientMessageChannel <- NewClientMessage(messageType, newClient.RoomID, newClient.ID, data)
		}
	})

	log.Printf("Server listened in port " + Port)

	if err := http.ListenAndServe(":"+Port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func (roomManager *RoomManager) RegisterRoom(roomId RoomId, room *Room) {
	roomManager.roomsMu.Lock()
	roomManager.rooms[roomId] = room
	roomManager.roomsMu.Unlock()
	log.Printf("룸 개수: %d", len(roomManager.rooms))
}

func (roomManager *RoomManager) setClient(client *Client) {
	roomManager.clientsMu.Lock()
	roomManager.clients[client.ID] = client
	roomManager.clientsMu.Unlock()
	room, exists := roomManager.rooms[client.RoomID]
	if !exists || room == nil {
		log.Printf("No Room: %s", client.RoomID)
		return
	}
	(*room).OnClient(string(client.ID))
}

func (roomManager *RoomManager) handleClientMessage(clientMessage *ClientMessage) {
	roomId := clientMessage.RoomId
	clientId := clientMessage.ClientId
	room, exists := roomManager.rooms[roomId]
	if !exists || room == nil {
		log.Printf("No Room: %s", roomId)
		return
	}
	(*room).OnMessage(clientMessage.Data, string(clientId))
}

// ######## Socket Interface 구현

func (roomManager *RoomManager) Send(data []byte, clientId string) error {
	return roomManager.clients[ClientId(clientId)].write(data)
}

func (roomManager *RoomManager) Broadcast(data []byte) {
	roomManager.clientsMu.Lock()
	for id, client := range roomManager.clients {
		if err := client.write(data); err != nil {
			log.Printf("Failed to send message to client %s: %v", id, err)
			client.close()
			delete(roomManager.clients, id)
		}
	}
	roomManager.clientsMu.Unlock()
}
