package game

import (
	"bytes"
	"encoding/json"
	"fmt"
	"jungle-royale/network"
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
	games                map[GameId]*Game
	gamesMu              sync.Mutex
	clientChannel        chan *Client
	clientMessageChannel chan *ClientMessage
	clientCloseChannel   chan *Client
}

// add game

func NewGameManager() *GameManager {
	socket := GameManager{
		make(map[GameId]*Game),
		sync.Mutex{},
		make(chan *Client, MaxClientCount),
		make(chan *ClientMessage, MaxClientCount),
		make(chan *Client, MaxClientCount),
	}
	return &socket
}

func (gameManager *GameManager) Listen() {

	go func() {
		for client := range gameManager.clientChannel {
			go gameManager.setClient(client)
		}
	}()

	go func() {
		for message := range gameManager.clientMessageChannel {
			go gameManager.handleClientMessage(message)
		}
	}()

	go func() {
		for message := range gameManager.clientCloseChannel {
			go gameManager.handleClientClose(message)
		}
	}()

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		// start := time.Now()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":200,"message":"Pong"}`)
		// elapsed := time.Since(start)
		// fmt.Printf("Request processed in %s\n", start.String())
	})

	http.HandleFunc("/api/create-game", func(w http.ResponseWriter, r *http.Request) {

		log.Printf("/api/create-game")

		w.Header().Set("Content-Type", "application/json")

		// HTTP 메서드 확인 (POST만 허용)
		if r.Method != http.MethodPost {
			http.Error(w, `{"status":405,"message":"Method Not Allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req network.GameServerNotificationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"status":400,"message":"Invalid JSON"}`, http.StatusBadRequest)
			return
		}

		gameManager.CreateGame(GameId(req.RoomID), req.MinPlayers, req.MaxPlayTime)

		response := `{"success":true,"message":"Game room created successfully"}`
		fmt.Fprintln(w, response)
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
		gameManager.clientChannel <- newClient

		log.Printf("Client %s connected", newClient.ID)

		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Client %s disconnected: %v", newClient.ID, err)
				gameManager.clientCloseChannel <- newClient
				// fmt.Printf("client disconnect\n")
				break
			}

			gameManager.clientMessageChannel <- NewClientMessage(messageType, newClient.GameID, newClient.ID, data)
		}
	})

	log.Printf("Server listened in port " + Port)

	if err := http.ListenAndServe(":"+Port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func (gameManager *GameManager) Test() {
	gameManager.sendStartMessage("test")
	gameManager.sendEndMessage("test")
}

func (gameManager *GameManager) sendStartMessage(gameId GameId) {
	url := "http://localhost:8080/api/game/start"

	// 요청 데이터 생성
	payload := network.StartMessageRequest{RoomID: string(gameId)}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		return
	}

	// HTTP POST 요청 보내기
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Start message response: %s\n", resp.Status)
}

func (gameManager *GameManager) sendEndMessage(gameId GameId) {
	url := "http://localhost:8080/api/game/end"

	// 요청 데이터 생성
	payload := network.EndMessageRequest{
		RoomID: string(gameId),
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		return
	}

	// HTTP POST 요청 보내기
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("End message response: %s\n", resp.Status)
}

func (gameManager *GameManager) CreateGame(
	roomId GameId,
	minPlayerNum int,
	playingTime int,
) {
	newRoom := NewGame(
		minPlayerNum,
		playingTime,
		func() {
			gameManager.sendStartMessage(roomId)
		},
		func() {
			gameManager.gamesMu.Lock()
			delete(gameManager.games, roomId)
			gameManager.gamesMu.Unlock()
			gameManager.sendEndMessage(roomId)
		},
	)
	newRoom.SetReadyStatus().StartGame() // 플레이어 수, 게임 시간
	gameManager.gamesMu.Lock()
	gameManager.games[roomId] = newRoom
	gameManager.gamesMu.Unlock()
	log.Printf("room: %d", len(gameManager.games))
}

func (gameManager *GameManager) setClient(client *Client) {
	room, exists := gameManager.games[client.GameID]
	if !exists || room == nil {
		log.Printf("No Room: %s", client.GameID)
		return
	}
	(*room).OnClient(client)
}

func (gameManager *GameManager) handleClientMessage(clientMessage *ClientMessage) {
	roomId := clientMessage.RoomId
	clientId := clientMessage.ClientId
	room, exists := gameManager.games[roomId]
	if !exists || room == nil {
		log.Printf("No Room: %s", roomId)
		return
	}
	(*room).OnMessage(clientMessage.Data, string(clientId))
}

func (gameManager *GameManager) handleClientClose(client *Client) {
	roomId := client.GameID
	room, exists := gameManager.games[roomId]
	if !exists || room == nil {
		log.Printf("No Room: %s", roomId)
		return
	}
	(*room).OnClose(client)
}
