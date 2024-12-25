package state

import (
	"jungle-royale/object"
	"math/rand"
	"sync"
)

type State struct {
	playersMu sync.Mutex
	Players   map[string]*object.Player
}

func NewState() *State {
	return &State{sync.Mutex{}, make(map[string]*object.Player)}
}

func (state *State) AddPlayer(id string) {
	newPlayer := object.NewPlayer(
		id,
		float32(rand.Intn(1000)),
		float32(rand.Intn(1000)),
	)
	state.playersMu.Lock()
	state.Players[id] = newPlayer
	state.playersMu.Unlock()
}

func (state *State) CalcState() {
	for idx := range state.Players {
		state.Players[idx].Move()
	}
}
