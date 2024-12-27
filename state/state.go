package state

import (
	"jungle-royale/chunk"
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object"
	"sync"

	"github.com/google/uuid"
)

type State struct {
	chunkNum  int
	chunkList [][]*chunk.Chunk
	MaxCoord  float32
	Players   sync.Map
	Bullets   sync.Map
}

func NewState() *State {
	return &State{
		Players: sync.Map{},
		Bullets: sync.Map{},
	}
}

func (state *State) SetState(chunkNum int) {
	state.chunkNum = chunkNum
	state.chunkList = make([][]*chunk.Chunk, chunkNum)
	for i := 0; i < chunkNum; i++ {
		state.chunkList[i] = make([]*chunk.Chunk, chunkNum)
		for j := 0; j < chunkNum; j++ {
			state.chunkList[i][j] = chunk.NewChunk()
		}
	}
	state.MaxCoord = float32(chunkNum * cons.CHUNK_LENGTH)
}

func (state *State) AddPlayer(id string, x float32, y float32) {
	newPlayer := object.NewPlayer(
		id,
		x,
		y,
	)
	state.Players.Store(id, newPlayer)
}

func (state *State) AddBullet(BulletCreateMessage *message.CreateBullet) {
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

func (state *State) CalcGameTickState() {

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

func (state *State) SecLoop() {

}
