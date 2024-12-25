package state

import (
	"jungle-royale/object"
	"math/rand"
	"sync"
)

type State struct {
	playersMu sync.Mutex
	players   []object.Player
}

func NewState() *State {
	return &State{sync.Mutex{}, []object.Player{}}
}

func (state *State) AddPlayer(id string) {
	state.players = append(state.players, *object.NewPlayer(id, float32(rand.Intn(1000)), float32(rand.Intn(1000))))
}

func (state *State) CalcLoop() {
	state.playersMu.Lock()
	for idx := range state.players {
		state.players[idx].Move()
	}
	state.playersMu.Unlock()
}
