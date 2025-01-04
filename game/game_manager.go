package game

import (
	"bytes"
	"encoding/json"
	"fmt"
	"jungle-royale/network"
	"jungle-royale/util"
	"log"
	"net"
	"net/http"
	"syscall"

	"github.com/gorilla/websocket"
)

const (
	MaxClientCount = 1000
	Port           = "8000"
)

type GameId string // room create 할 때 받음

type GameManager struct {
	games                *util.Map[GameId, *Game]
	clientChannel        chan *Client
	clientMessageChannel chan *ClientMessage
	clientCloseChannel   chan *Client
	debug                bool // production, development 환경 체크
}

// add game

func NewGameManager(
	debug bool,
) *GameManager {
	socket := GameManager{
		util.NewSyncMap[GameId, *Game](),
		make(chan *Client, MaxClientCount),
		make(chan *ClientMessage, MaxClientCount),
		make(chan *Client, MaxClientCount),
		debug,
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
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":200,"message":"Pong"}`)
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

		gameManager.CreateGame(GameId(req.GameID), req.MinPlayers, req.MaxPlayTime)

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
		defer conn.Close()

		// Nagle 알고리즘 끄기
		tcpConn := gameManager.getTCPConn(conn)
		if tcpConn != nil {
			if err := gameManager.setNoDelay(tcpConn); err != nil {
				log.Println("Failed to disable Nagle:", err)
			} else {
				log.Println("Nagle disabled")
			}
		}

		gameId := r.URL.Query().Get("roomId")
		if gameId == "" {
			http.Error(w, "Missing gameId query parameter", http.StatusBadRequest)
			return
		}

		serverClientId := r.URL.Query().Get("clientId")
		if serverClientId == "" {
			http.Error(w, "Missing serverClientId query parameter", http.StatusBadRequest)
			return
		}

		log.Println("new client", gameId, serverClientId)
		newClient := NewClient(GameId(gameId), serverClientId, conn)
		gameManager.clientChannel <- newClient

		log.Printf("Client %s connected", newClient.ID)

		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(
					err,
					websocket.CloseGoingAway,
					websocket.CloseAbnormalClosure,
				) {
					log.Printf("Client %s disconnected unexpectedly: %v", newClient.ID, err)
				} else {
					log.Printf("Client %s disconnected: %v", newClient.ID, err)
				}
				gameManager.clientCloseChannel <- newClient
				break
			}

			gameManager.clientMessageChannel <- NewClientMessage(messageType, newClient.GameID, newClient.ID, data)
		}
	})

	log.Printf("Server listened in port " + Port)

	server := &http.Server{Addr: ":" + Port}
	if err := server.ListenAndServe(); err != nil {

		gameManager.games.Range(func(gi GameId, g *Game) bool {
			g.clients.Range(func(ci ClientId, c *Client) bool {
				c.close()
				return true
			})
			return true
		})

		log.Printf("Server failed: %v", err)
	}
}

func (gameManager *GameManager) sendStartMessage(gameId GameId) {
	url := "http://wep-api.eternalsnowman.com:8080"
	if gameManager.debug {
		url = "http://localhost:8080"
	}
	url += "/api/game/start"

	// 요청 데이터 생성
	payload := network.StartMessageRequest{GameID: string(gameId)}
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
	url := "http://wep-api.eternalsnowman.com:8080"
	if gameManager.debug {
		url = "http://localhost:8080"
	}
	url += "/api/game/end"

	// 요청 데이터 생성
	payload := network.EndMessageRequest{
		GameID: string(gameId),
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
	gameId GameId,
	minPlayerNum int,
	playingTime int,
) {
	newGame := NewGame(
		minPlayerNum,
		playingTime,
		func() { // 게임 시작 (대기방에서 시작화면으로)
			gameManager.handleGameStart(gameId)
		},
		func() { // 게임 종료
			gameManager.handleGameEnd(gameId)
		},
	)
	newGame.SetReadyStatus().StartGame() // 플레이어 수, 게임 시간
	gameManager.games.Store(gameId, newGame)
	log.Printf("room count: %d", gameManager.games.Length())
}

func (gameManager *GameManager) handleGameStart(gameId GameId) {
	if gameManager.debug {
		return
	}
	gameManager.sendStartMessage(gameId)
}

func (gameManager *GameManager) handleGameEnd(gameId GameId) {
	if gameManager.debug {
		return
	} else {
		gameManager.games.Delete(gameId)
		gameManager.sendEndMessage(gameId)
	}
}

func (gameManager *GameManager) setClient(client *Client) {
	room, exists := gameManager.games.Get(client.GameID)
	if !exists || room == nil {
		log.Printf("No Room client: %s", client.GameID)
		return
	}
	(*room).OnClient(client)
}

func (gameManager *GameManager) handleClientMessage(clientMessage *ClientMessage) {
	gameId := clientMessage.GameId
	clientId := clientMessage.ClientId
	room, exists := gameManager.games.Get(gameId)
	if !exists || room == nil {
		log.Printf("No Room clinet message: %s", gameId)
		return
	}
	(*room).OnMessage(clientMessage.Data, string(clientId))
}

func (gameManager *GameManager) handleClientClose(client *Client) {
	gameId := client.GameID
	room, exists := gameManager.games.Get(gameId)
	if !exists || room == nil {
		log.Printf("No Room close: %s", gameId)
		return
	}
	(*room).OnClose(client)
}

// TCP 연결 가져오기
func (gameManager *GameManager) getTCPConn(conn *websocket.Conn) *net.TCPConn {
	rawConn := conn.UnderlyingConn() // WebSocket의 net.Conn 가져오기
	if tcpConn, ok := rawConn.(*net.TCPConn); ok {
		return tcpConn
	}
	return nil
}

// Nagle 알고리즘 끄기
func (gameManager *GameManager) setNoDelay(conn *net.TCPConn) error {
	rawConn, err := conn.SyscallConn()
	if err != nil {
		return err
	}

	return rawConn.Control(func(fd uintptr) {
		syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1)
	})
}
