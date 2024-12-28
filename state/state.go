package state

import (
	"jungle-royale/chunk"
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object"

	"github.com/google/uuid"
)

type State struct {
	chunkNum   int
	chunkList  [][]*chunk.Chunk
	MaxCoord   float32
	ObjectList *object.SyncMapList
}

func NewState() *State {
	return &State{
		ObjectList: object.NewMoverSyncMapList(),
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
	state.ObjectList.GetPlayers().Store(id, newPlayer)
}

func (state *State) AddBullet(BulletCreateMessage *message.CreateBullet) {
	bulletId := uuid.New().String()
	newBullet := object.NewBullet(
		bulletId,
		BulletCreateMessage.PlayerId,
		BulletCreateMessage.StartX,
		BulletCreateMessage.StartY,
		float64(BulletCreateMessage.Angle),
	)
	state.ObjectList.GetBullets().Store(bulletId, newBullet)
}
