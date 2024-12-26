package state

import (
	"jungle-royale/message"
	"jungle-royale/object"
	"math/rand"
	"sync"

	"github.com/google/uuid"
)

type State struct {
	Players sync.Map
	Bullets sync.Map
}

func NewState() *State {
	return &State{
		Players: sync.Map{},
		Bullets: sync.Map{},
	}
}

func (state *State) AddPlayer(id string) {
	newPlayer := object.NewPlayer(
		id,
		float32(rand.Intn(1000)),
		float32(rand.Intn(1000)),
	)
	state.Players.Store(id, newPlayer)
}

func (state *State) AddBullet(BulletCreateMessage *message.BulletCreate) {
	bulletId := uuid.New().String()
	newBullet := object.NewBullet(
		bulletId,
		BulletCreateMessage.PlayerId,
		BulletCreateMessage.StartX,
		BulletCreateMessage.StartY,
		float64(BulletCreateMessage.Angle),
		state.GetPlayers(),
	)
	state.Bullets.Store(bulletId, newBullet)
}

func (state *State) GetPlayers() map[string]*object.Player {
	players := make(map[string]*object.Player)
	state.Players.Range(func(key, value any) bool {
		players[key.(string)] = value.(*object.Player)
		return true
	})
	return players
}

func (state *State) CalcState() {

	state.Players.Range(func(key, value any) bool {
		playerId := key.(string)
		player := value.(*object.Player)
		if !player.IsValid() {
			state.Players.Delete(playerId)
			return true
		}
		player.Move()
		return true
	})

	state.Bullets.Range(func(key, value any) bool {
		bulletId := key.(string)
		bullet := value.(*object.Bullet)
		isValid := bullet.Move()
		if !isValid {
			state.Bullets.Delete(bulletId)
		}
		return true
	})
}
