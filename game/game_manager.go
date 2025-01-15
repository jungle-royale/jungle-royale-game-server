package game

import (
	"bytes"
	"encoding/json"
	"fmt"
	"jungle-royale/cons"
	"jungle-royale/network"
	"jungle-royale/serverlog"
	"jungle-royale/util"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"syscall"

	"github.com/gorilla/websocket"
	"golang.org/x/net/http2"
)

const Port = "8000"

type GameId string // room create 할 때 받음

type GameManager struct {
	games                *util.Map[GameId, int]
	gameRooms            []*Game // GameId: idx
	gameRoomLock         sync.Mutex
	emptyGameIdx         int
	clientChannel        chan *Client
	clientMessageChannel chan *ClientMessage
	clientCloseChannel   chan *Client
	debug                bool // production, development 환경 체크
	debugClientCount     int  // debug 전용
	gameManagerLogger    *serverlog.GameManagerLogger
	ClientIdAllocator    *ClientIdAllocator
}

func NewGameManager(
	debug bool,
) *GameManager {
	socket := GameManager{
		util.NewSyncMap[GameId, int](),
		make([]*Game, cons.MaxGameNum),
		sync.Mutex{},
		0,
		make(chan *Client, cons.MaxClientCount),
		make(chan *ClientMessage, cons.MaxClientCount),
		make(chan *Client, cons.MaxClientCount),
		debug,
		0,
		serverlog.NewGameManagerLogger(),
		NewClientIdAllocator(),
	}
	for i := 0; i < cons.MaxGameNum; i++ {
		socket.gameRooms[i] = NewGame(debug)
	}
	socket.gameManagerLogger.Log("New GameManager")
	return &socket
}

func (gameManager *GameManager) Listen() {

	go func() {
		for client := range gameManager.clientChannel {
			if client.isObserver {
				go gameManager.setObserver(client)
			} else {
				go gameManager.setClient(client)
			}
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
		// CORS 헤더 추가
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// OPTIONS 요청에 대한 처리
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		fmt.Fprintf(w, `{"status":200,"message":"Pong"}`)
	})

	http.HandleFunc("/api/create-game", func(w http.ResponseWriter, r *http.Request) {

		// log.Printf("/api/create-game")
		gameManager.gameManagerLogger.Log("request /api/create-game")

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

		gameId := GameId(req.GameID)

		_, exists := gameManager.games.Get(gameId)
		if exists {
			// log.Printf("이미 존재하는 게임: %s", gameId)
			gameManager.gameManagerLogger.Log("이미 존재하는 게임: " + string(gameId))
			response := `{"success":false,"message":"Game room created successfully"}`
			fmt.Fprintln(w, response)
		} else {
			if !gameManager.SetNewGame(gameId, req.MinPlayers, req.MaxPlayTime) {
				// log.Printf("모든 룸에서 게임 진행 중")
				gameManager.gameManagerLogger.Log("모든 룸에서 게임 진행 중")
			} else {
				response := `{"success":true,"message":"Game room created successfully"}`
				fmt.Fprintln(w, response)
			}
		}
	})

	http.HandleFunc("/observer", func(w http.ResponseWriter, r *http.Request) {

		gameManager.gameManagerLogger.Log("request " + r.URL.String())

		var upgrader = websocket.Upgrader{
			ReadBufferSize:  8192, // 읽기 버퍼 크기
			WriteBufferSize: 8192, // 쓰기 버퍼 크기
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to upgrade to WebSocket: %v", err)
			return
		}
		defer conn.Close()

		tcpConn := gameManager.getTCPConn(conn)
		if tcpConn != nil {
			if err := gameManager.setNoDelay(tcpConn); err != nil {
				log.Println("Failed to disable Nagle:", err)
			}
		}

		var gameId string

		if gameManager.debug {
			gameId = "test"
			gameManager.debugClientCount += 1
		} else {
			gameId = r.URL.Query().Get("roomId")
			if gameId == "" {
				http.Error(w, "Missing gameId query parameter", http.StatusBadRequest)
				return
			}
		}

		// log.Println("new client", gameId, serverClientId)
		gameManager.gameManagerLogger.Log("new observer")

		newClient := NewClient(GameId(gameId), "", "", conn, true, gameManager.ClientIdAllocator.AllocateClientId())
		gameManager.clientChannel <- newClient

		log.Printf("Client %d connected", newClient.ID)

		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(
					err,
					websocket.CloseGoingAway,
					websocket.CloseAbnormalClosure,
				) {
					log.Printf("Client %d disconnected unexpectedly: %v", newClient.ID, err)
				} else {
					log.Printf("Client %d disconnected: %v", newClient.ID, err)
				}
				gameManager.clientCloseChannel <- newClient
				break
			}

			gameManager.clientMessageChannel <- NewClientMessage(messageType, newClient.GameID, newClient.ID, data)
		}
	})

	http.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) {

		// log.Printf("%s", r.URL.String())
		gameManager.gameManagerLogger.Log("request " + r.URL.String())

		var upgrader = websocket.Upgrader{
			ReadBufferSize:  8192, // 읽기 버퍼 크기
			WriteBufferSize: 8192, // 쓰기 버퍼 크기
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

		tcpConn := gameManager.getTCPConn(conn)
		if tcpConn != nil {
			if err := gameManager.setNoDelay(tcpConn); err != nil {
				log.Println("Failed to disable Nagle:", err)
			}
		}

		var gameId string
		var serverClientId string
		var userName string

		if gameManager.debug {
			gameId = "test"
			serverClientId = strconv.Itoa(gameManager.debugClientCount)
			gameManager.debugClientCount += 1
		} else {
			gameId = r.URL.Query().Get("roomId")
			if gameId == "" {
				http.Error(w, "Missing gameId query parameter", http.StatusBadRequest)
				return
			}
			serverClientId = r.URL.Query().Get("clientId")
			if serverClientId == "" {
				http.Error(w, "Missing serverClientId query parameter", http.StatusBadRequest)
				return
			}
			userName = r.URL.Query().Get("username")
			if userName == "" {
				http.Error(w, "Missing userName query parameter", http.StatusBadRequest)
				return
			}
		}

		// log.Println("new client", gameId, serverClientId)
		gameManager.gameManagerLogger.Log("new client " + serverClientId)

		newClient := NewClient(GameId(gameId), serverClientId, userName, conn, false, gameManager.ClientIdAllocator.AllocateClientId())
		gameManager.clientChannel <- newClient

		log.Printf("Client %d connected", newClient.ID)

		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(
					err,
					websocket.CloseGoingAway,
					websocket.CloseAbnormalClosure,
				) {
					log.Printf("Client %d disconnected unexpectedly: %v", newClient.ID, err)
				} else {
					log.Printf("Client %d disconnected: %v", newClient.ID, err)
				}
				gameManager.clientCloseChannel <- newClient
				break
			}

			gameManager.clientMessageChannel <- NewClientMessage(messageType, newClient.GameID, newClient.ID, data)
		}
	})

	gameManager.SetServerManager()

	// log.Printf("Server listened in port " + Port)
	gameManager.gameManagerLogger.Log("Server listened in port " + Port)

	server := &http.Server{Addr: ":" + Port}
	http2.ConfigureServer(server, &http2.Server{})
	if err := server.ListenAndServe(); err != nil {

		gameManager.games.Range(func(gi GameId, i int) bool {
			g := gameManager.gameRooms[i]
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
	url := "http://web-api.eternalsnowman.com:8080"
	if gameManager.debug {
		url = "http://localhost:8080"
	}
	url += "/api/game/start"

	// 요청 데이터 생성
	gameIdx, _ := gameManager.games.Get(gameId)
	game := gameManager.gameRooms[*gameIdx]
	clientIds := make([]string, 0)
	game.clients.Range(func(ci ClientId, c *Client) bool {
		clientIds = append(clientIds, c.serverClientId)
		return true
	})
	payload := network.StartMessageRequest{
		GameID:    string(gameId),
		ClientIds: clientIds,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		// fmt.Printf("Error encoding JSON: %v\n", err)
		gameManager.gameManagerLogger.Log("Error encoding JSON: " + err.Error())
		return
	}

	// HTTP POST 요청 보내기
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		// fmt.Printf("Error sending request: %v\n", err)
		gameManager.gameManagerLogger.Log("Error sending request: " + err.Error())
		return
	}
	defer resp.Body.Close()

	// fmt.Printf("Start message response: %s\n", resp.Status)
	gameManager.gameManagerLogger.Log("Start message response: " + resp.Status)
}

func (gameManager *GameManager) sendEndMessage(gameId GameId) {
	url := "http://web-api.eternalsnowman.com:8080"
	// url := "http://172.16.155.128:8080"
	if gameManager.debug {
		url = "http://localhost:8080"
	}
	url += "/api/game/end"

	// 요청 데이터 생성
	gameIdx, ok := gameManager.games.Get(gameId)
	if !ok {
		// log.Printf("no room: %s", gameId)
		gameManager.gameManagerLogger.Log("no room: " + string(gameId))
		return
	}

	game := gameManager.gameRooms[*gameIdx]
	payload := network.EndMessageRequest{
		GameID:  string(gameId),
		GameLog: game.gameRecorder.ReturnList(),
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		// fmt.Printf("Error encoding JSON: %v\n", err)
		gameManager.gameManagerLogger.Log("Error encoding JSON: " + err.Error())
		return
	}

	// HTTP POST 요청 보내기
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		// fmt.Printf("Error sending request: %v\n", err)
		gameManager.gameManagerLogger.Log("Error sending request: " + err.Error())
		return
	}
	defer resp.Body.Close()

	// fmt.Printf("End message response: %s\n", resp.Status)
	gameManager.gameManagerLogger.Log("End message response: " + resp.Status)
}

func (gameManager *GameManager) sendPlayerLeaveMessage(gameId GameId, client *Client) {
	url := "http://web-api.eternalsnowman.com:8080"
	// url := "http://172.16.155.128:8080"
	if gameManager.debug {
		url = "http://localhost:8080"
	}
	url += "/api/game/leave"

	// 요청 데이터 생성
	msg := network.PlayerLeaveMessageRequest{
		GameID:   string(gameId),
		ClientID: client.serverClientId,
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		// fmt.Printf("Error encoding JSON: %v\n", err)
		gameManager.gameManagerLogger.Log("Error encoding JSON: " + err.Error())
		return
	}

	// HTTP POST 요청 보내기
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		// fmt.Printf("Error sending request: %v\n", err)
		gameManager.gameManagerLogger.Log("Error sending request: " + err.Error())
		return
	}
	defer resp.Body.Close()

	// fmt.Printf("leave message response: %s\n", resp.Status)
	gameManager.gameManagerLogger.Log("End message response: " + resp.Status)
}

func (gameManager *GameManager) SetNewGame(
	gameId GameId,
	minPlayerNum int,
	playingTime int,
) bool {
	gameManager.gameRoomLock.Lock()
	newGame := gameManager.gameRooms[gameManager.emptyGameIdx]
	gameManager.games.Store(gameId, gameManager.emptyGameIdx)
	newGame.SetGame(
		minPlayerNum,
		playingTime,
		func() { // 게임 시작 (대기방에서 시작화면으로)
			gameManager.handleGameStart(gameId)
		},
		func(client *Client) { // 플레이어 대기방 떠남
			gameManager.handlePlayerLeave(gameId, client)
		},
		func() { // 게임 종료
			gameManager.handleGameEnd(gameId)
		},
		serverlog.NewGameLogger(gameManager.emptyGameIdx, string(gameId)),
	)
	go newGame.SetReadyStatus().StartGame()
	last := gameManager.emptyGameIdx
	for {
		gameManager.emptyGameIdx = (gameManager.emptyGameIdx + 1) % cons.MaxGameNum
		if gameManager.emptyGameIdx == last {
			return false
		}
		if !gameManager.gameRooms[gameManager.emptyGameIdx].IsPlaying() {
			break
		}
	}
	gameManager.gameRoomLock.Unlock()
	return true
}

func (gameManager *GameManager) handleGameStart(gameId GameId) {
	if gameManager.debug {
		return
	}
	gameManager.sendStartMessage(gameId)
	log.Printf("Start Game: %s", gameId)
}

func (gameManager *GameManager) handleGameEnd(gameId GameId) {
	if gameManager.debug {
		return
	} else {
		gameIdx, _ := gameManager.games.Get(gameId)
		gameManager.sendEndMessage(gameId)
		gameManager.games.Delete(gameId)
		gameManager.gameRooms[*gameIdx] = NewGame(gameManager.debug)
	}
}

func (gameManager *GameManager) handlePlayerLeave(gameId GameId, client *Client) {
	if gameManager.debug {
		return
	} else {
		gameManager.sendPlayerLeaveMessage(gameId, client)
	}
}

func (gameManager *GameManager) setClient(client *Client) {
	if idx, exists := gameManager.games.Get(client.GameID); exists {
		room := gameManager.gameRooms[*idx]
		if room.IsPlaying() {
			room.OnClient(client)
			return
		}
	}

	log.Printf("New Client / No Room: client is.. %s, %d", client.GameID, gameManager.games.Length())
	client.close()
}

func (gameManager *GameManager) setObserver(client *Client) {
	if idx, exists := gameManager.games.Get(client.GameID); exists {
		room := gameManager.gameRooms[*idx]
		if room.IsPlaying() {
			room.OnObserver(client)
			return
		}
	}
}

func (gameManager *GameManager) handleClientMessage(clientMessage *ClientMessage) {
	gameId := clientMessage.GameId
	clientId := clientMessage.ClientId
	if idx, exists := gameManager.games.Get(gameId); exists {
		room := gameManager.gameRooms[*idx]
		if room.IsPlaying() {
			room.OnMessage(clientMessage.Data, int(clientId))
			return
		}
	}

	log.Printf("No Room: clinet message is.. %s", gameId)
}

func (gameManager *GameManager) handleClientClose(client *Client) {
	gameId := client.GameID
	if idx, exists := gameManager.games.Get(gameId); exists {
		room := gameManager.gameRooms[*idx]
		if room.IsPlaying() {
			room.OnClose(client)
			return
		}
	}

	log.Printf("No Room close: %s", gameId)
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
