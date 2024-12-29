package state

import (
	"jungle-royale/chunk"
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object"
	"jungle-royale/util"

	"github.com/google/uuid"
)

type State struct {
	chunkNum   int
	chunkList  [][]*chunk.Chunk
	Players    *util.Map[string, *object.Player]
	Bullets    *util.Map[string, *object.Bullet]
	HealPacks  *util.Map[string, *object.HealPack]
	MagicItems *util.Map[string, *object.Magic]
	MaxCoord   float32
}

func NewState() *State {
	return &State{
		Players:    util.NewSyncMap[string, *object.Player](),
		Bullets:    util.NewSyncMap[string, *object.Bullet](),
		HealPacks:  util.NewSyncMap[string, *object.HealPack](),
		MagicItems: util.NewSyncMap[string, *object.Magic](),
	}
}

func (state *State) ConfigureState(chunkNum int) {
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
	if player, exists := state.Players.Get(BulletCreateMessage.PlayerId); exists {
		newBullet := object.NewBullet(
			bulletId,
			BulletCreateMessage.PlayerId,
			(*player).MagicType,
			BulletCreateMessage.StartX,
			BulletCreateMessage.StartY,
			float64(BulletCreateMessage.Angle),
		)
		state.Bullets.Store(bulletId, newBullet)
	}
}

func (state *State) ChangeDirection(clientId string, msg *message.ChangeDir) {
	if player, exists := state.Players.Get(clientId); exists {
		(*player).DirChange(float64(msg.GetAngle()), msg.IsMoved)
	}
}

func (state *State) CreateBullet(clientId string, msg *message.CreateBullet) {
	state.AddBullet(msg)
}

func (state *State) DoDash(clientId string, msg *message.DoDash) {
	if player, exists := state.Players.Get(clientId); exists {
		(*player).DoDash()
	}
}
