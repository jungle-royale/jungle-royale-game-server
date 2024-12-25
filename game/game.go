package game

import (
	"jungle-royale/message"
	"jungle-royale/socket"
	"jungle-royale/state"
	"log"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/gorilla/websocket"
)

type Game struct {
	state     *state.State
	clients   []*socket.Client
	clientsMu sync.Mutex
}

func NewGame() *Game {
	GameContext := &Game{state.NewState(), []*socket.Client{}, sync.Mutex{}} // generate game
	go (*GameContext).CalcLoop()                                             // start main loop
	go (*GameContext).BroadcastLoop()                                        // broadcast to client

	// connection event handler
	go func() {
		for client := range socket.ClientChannel {
			go GameContext.SetPlayer(client)
		}
	}()

	// message event handler
	go func() {
		for message := range socket.MessageChannel {
			go GameContext.HandleMessage(message)
		}
	}()

	return GameContext
}

func (game *Game) SetPlayer(client *socket.Client) {

	game.clientsMu.Lock()
	game.clients = append(game.clients, client)
	game.clientsMu.Unlock()

	game.state.AddPlayer(client.ID)

	// send GameInit message
	gameInit := &message.GameInit{Id: client.ID}
	data, err := proto.Marshal(&message.Wrapper{
		MessageType: &message.Wrapper_Gameinit{
			Gameinit: gameInit,
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

func (game *Game) HandleMessage(clientMessage *socket.ClientMessage) {

	if clientMessage.MessageType == websocket.BinaryMessage {

		// Protobuf message decode
		var wrapper message.Wrapper
		if err := proto.Unmarshal(clientMessage.Data, &wrapper); err != nil {
			log.Printf("Failed to unmarshal message from client %s: %v", clientMessage.ID, err)
			return
		}

		// change message
		if change := wrapper.GetChange(); change != nil {
			if player, exists := game.state.Players[clientMessage.ID]; exists {
				go player.DirChange(change.Dx, change.Dy)
			}
		}
	}
}

func (game *Game) CalcLoop() {
	ticker := time.NewTicker(16 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C { // calculation loop
		game.state.CalcState()
	}
}

func (game *Game) BroadcastLoop() {
	ticker := time.NewTicker(16 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C { // broadcast loop

		playerList := make([]*message.Player, 0, len(game.state.Players))
		for _, player := range game.state.Players {
			playerList = append(playerList, player.MakeSendingPlayerData())
		}

		gameState := &message.GameState{
			Players: playerList,
		}

		data, err := proto.Marshal(&message.Wrapper{
			MessageType: &message.Wrapper_State{
				State: gameState,
			},
		})
		if err != nil {
			log.Printf("Falied to marshal GameState: %v", err)
			return
		}

		for _, client := range game.clients {
			if err := client.Conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
				log.Printf("Failed to send message to client %s: %v", client.ID, err)
				client.Conn.Close()
			}
		}
	}
}
