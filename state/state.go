package state

import (
	"jungle-royale/chunk"
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object"

	"github.com/google/uuid"
)

type State struct {
	chunkNum  int
	chunkList [][]*chunk.Chunk
	MaxCoord  float32
	MoverList *object.MoverSyncMapList
}

func NewState() *State {
	return &State{
		MoverList: object.NewMoverSyncMapList(),
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
	state.MoverList.GetPlayers().Store(id, newPlayer)
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
	state.MoverList.GetBullets().Store(bulletId, newBullet)
}

func (state *State) GetPlayers() map[string]*object.Player {
	players := make(map[string]*object.Player)
	state.MoverList.GetPlayers().Range(func(key, value any) bool {
		players[key.(string)] = value.(*object.Player)
		return true
	})
	return players
}

func (state *State) CalcGameTickState() {

	// player
	state.MoverList.GetPlayers().Range(func(key, value any) bool {
		playerId := key.(string)
		player := value.(*object.Player)
		if !player.IsValid() {
			state.MoverList.GetPlayers().Delete(playerId)
			return true
		}
		player.CalcGameTick()
		return true
	})

	// bullet
	state.MoverList.GetBullets().Range(func(key, value any) bool {
		bulletId := key.(string)
		bullet := value.(*object.Bullet)
		bullet.CalcGameTick()
		if !bullet.IsValid() {
			state.MoverList.GetBullets().Delete(bulletId)
		}
		return true
	})
}

func (state *State) SecLoop() {

}
