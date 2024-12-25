package game

import (
	"fmt"
	"jungle-royale/socket"
	"jungle-royale/state"
	"time"
)

type Game struct {
	state   *state.State
	clients []*socket.Client
}

func NewGame() *Game {
	GameContext := &Game{state.NewState(), []*socket.Client{}} // generate game
	go (*GameContext).CalcLoop()                               // start main loop
	go (*GameContext).BroadcastLoop()                          // broadcast to client

	// connection event handler
	go func() {
		for conn := range socket.ConnectionChannel {
			id := generateSessionID()
			GameContext.clients = append(GameContext.clients, socket.NewClient(conn, id))
		}
	}()

	return GameContext
}

func generateSessionID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func (game *Game) CalcLoop() {
	ticker := time.NewTicker(16 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C { // calculation loop

	}
}

func (game *Game) BroadcastLoop() {
	ticker := time.NewTicker(16 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C { // broadcast loop

	}
}
