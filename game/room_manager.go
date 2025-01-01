package game

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	MaxClientCount = 1000
	Port           = "8000"
)

type GameId string // room create 할 때 받음

type GameManager struct {
	rooms                map[GameId]*Room
	roomsMu              sync.Mutex
	clientChannel        chan *Client
	clientMessageChannel chan *ClientMessage
}

// add game

func NewGameManager() *GameManager {
	socket := GameManager{
		make(map[GameId]*Room),
		sync.Mutex{},
		make(chan *Client, MaxClientCount),
		make(chan *ClientMessage, MaxClientCount),
	}
	return &socket
}

func (socket *GameManager) Listen() {

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

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		// start := time.Now()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":200,"message":"Pong"}`)
		// elapsed := time.Since(start)
		// fmt.Printf("Request processed in %s\n", start.String())
	})

	http.HandleFunc("/game/start", func(w http.ResponseWriter, r *http.Request) {
		// start := time.Now()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":200,"message":"Pong"}`)
		// elapsed := time.Since(start)
		// fmt.Printf("Request processed in %s\n", start.String())
	})

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

		newClient := NewClient(GameId(roomId), conn)
		socket.clientChannel <- newClient

		log.Printf("Client %s connected", newClient.ID)

		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Client %s disconnected: %v", newClient.ID, err)
				// fmt.Printf("client disconnect\n")
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

func (gameManager *GameManager) CreateRoom(
	roomId GameId,
	minPlayerNum int,
	playingTime int,
) {
	// gameManager := game.NewgameManager()
	var newRoom Room = NewGame(minPlayerNum, playingTime).SetReadyStatus().StartGame() // 플레이어 수, 게임 시간
	gameManager.roomsMu.Lock()
	gameManager.rooms[roomId] = &newRoom
	gameManager.roomsMu.Unlock()
	log.Printf("room: %d", len(gameManager.rooms))
}

func (gameManager *GameManager) setClient(client *Client) {
	room, exists := gameManager.rooms[client.RoomID]
	if !exists || room == nil {
		log.Printf("No Room: %s", client.RoomID)
		return
	}
	(*room).OnClient(client)
}

func (gameManager *GameManager) handleClientMessage(clientMessage *ClientMessage) {
	roomId := clientMessage.RoomId
	clientId := clientMessage.ClientId
	room, exists := gameManager.rooms[roomId]
	if !exists || room == nil {
		log.Printf("No Room: %s", roomId)
		return
	}
	(*room).OnMessage(clientMessage.Data, string(clientId))
}
