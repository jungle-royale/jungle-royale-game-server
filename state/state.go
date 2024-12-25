package state

import (
	"jungle-royale/message"
	"jungle-royale/object"
	"math/rand"
	"sync"

	"github.com/google/uuid"
)

type State struct {
	playersMu sync.Mutex
	Players   map[string]*object.Player
	bulletsMu sync.Mutex
	Bullets   map[string]*object.Bullet
}

func NewState() *State {
	return &State{
		sync.Mutex{},
		make(map[string]*object.Player),
		sync.Mutex{},
		make(map[string]*object.Bullet),
	}
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

func (state *State) AddBullet(BulletCreateMessage *message.BulletCreate) {
	bulletId := uuid.New().String()
	newBullet := object.NewBullet(
		bulletId,
		BulletCreateMessage.PlayerId,
		BulletCreateMessage.StartX,
		BulletCreateMessage.StartY,
		float64(BulletCreateMessage.Angle),
		state.Players,
	)
	state.bulletsMu.Lock()
	state.Bullets[bulletId] = newBullet
	state.bulletsMu.Unlock()
}

func (state *State) CalcState() {
	for id := range state.Players {
		state.Players[id].Move()
	}
	for id := range state.Bullets {
		isValie := state.Bullets[id].Move()
		if !isValie {
			delete(state.Bullets, id)
		}
	}
}
