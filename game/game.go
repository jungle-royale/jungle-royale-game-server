package game

import (
	"jungle-royale/message"
	"jungle-royale/object"
	"jungle-royale/state"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

const (
	MaxClientCount        = 100
	CalcLoopInterval      = 16
	BroadCastLoopInterval = 16
)

type Game struct {
	state                *state.State
	clients              map[string]*GameClient
	clientChannel        chan *GameClient
	clientMessageChannel chan *GameClientMessage
	clientsMu            sync.Mutex
}

func NewGame() *Game {
	return &Game{
		state.NewState(),
		make(map[string]*GameClient),
		make(chan *GameClient, MaxClientCount),
		make(chan *GameClientMessage, MaxClientCount),
		sync.Mutex{},
	}
}

func (game *Game) StartGame() {
	go game.CalcLoop()      // start main loop
	go game.BroadcastLoop() // broadcast to client

	// connection event handler
	go func() {
		for client := range game.clientChannel {
			go game.SetPlayer(client)
		}
	}()

	// message event handler
	go func() {
		for message := range game.clientMessageChannel {
			go game.HandleMessage(message)
		}
	}()

	game.InitSocket(func(w http.ResponseWriter, r *http.Request) {
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

		newClient := NewGameClient(conn)

		game.clientChannel <- newClient

		log.Printf("Client %s connected", newClient.ID)

		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Client %s disconnected: %v", newClient.ID, err)
				break
			}

			game.clientMessageChannel <- NewGameClientMessage(messageType, newClient.ID, data)
		}
	})
}

func (game *Game) InitSocket(handleWebSocket func(w http.ResponseWriter, r *http.Request)) {
	http.HandleFunc("/ws", handleWebSocket)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
	log.Printf("Server listened in port 8000")
}

func (game *Game) SetPlayer(client *GameClient) {
	game.clientsMu.Lock()
	game.clients[client.ID] = client
	game.clientsMu.Unlock()

	game.state.AddPlayer(client.ID)

	// send GameInit message
	gameInit := &message.GameInit{Id: client.ID}
	data, err := proto.Marshal(&message.Wrapper{
		MessageType: &message.Wrapper_GameInit{
			GameInit: gameInit,
		},
	})
	if err != nil {
		log.Printf("Failed to marshal GameInit: %v", err)
		return
	}
	if err := client.Conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		log.Printf("Failed to send GameInit message to client %s: %v", client.ID, err)
	}
}

func (game *Game) HandleMessage(clientMessage *GameClientMessage) {
	if clientMessage.MessageType == websocket.BinaryMessage {
		// Protobuf message decode
		var wrapper message.Wrapper
		if err := proto.Unmarshal(clientMessage.Data, &wrapper); err != nil {
			log.Printf("Failed to unmarshal message from client %s: %v", clientMessage.ID, err)
			return
		}

		// dirChange message
		if dirChange := wrapper.GetDirChange(); dirChange != nil {
			if value, exists := game.state.Players.Load(clientMessage.ID); exists {
				player := value.(*object.Player)
				go player.DirChange(float64(dirChange.GetAngle()), dirChange.IsMoved)
			}
		}

		// bulletCreate message
		if bulletCreate := wrapper.GetBulletCreate(); bulletCreate != nil {
			game.state.AddBullet(bulletCreate)
		}
	}
}

func (game *Game) CalcLoop() {
	ticker := time.NewTicker(CalcLoopInterval * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C { // calculation loop
		game.state.CalcState()
	}
}

func (game *Game) BroadcastLoop() {
	ticker := time.NewTicker(BroadCastLoopInterval * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C { // broadcast loop
		playerList := make([]*message.Player, 0)
		game.state.Players.Range(func(key, value any) bool {
			player := value.(*object.Player)
			playerList = append(playerList, player.MakeSendingData())
			return true
		})

		bulletList := make([]*message.BulletState, 0)
		game.state.Bullets.Range(func(key, value any) bool {
			bullet := value.(*object.Bullet)
			bulletList = append(bulletList, bullet.MakeSendingData())
			return true
		})

		gameState := &message.GameState{
			Players:     playerList,
			BulletState: bulletList,
		}

		data, err := proto.Marshal(&message.Wrapper{
			MessageType: &message.Wrapper_State{
				State: gameState,
			},
		})
		if err != nil {
			log.Printf("Failed to marshal GameState: %v", err)
			return
		}

		game.clientsMu.Lock()
		for id, client := range game.clients {
			if err := client.Conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
				log.Printf("Failed to send message to client %s: %v", client.ID, err)
				client.Conn.Close()
				delete(game.clients, id)
			}
		}
		game.clientsMu.Unlock()
	}
}
