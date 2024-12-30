package state

import (
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object"
	"jungle-royale/util"

	"github.com/google/uuid"
)

type State struct {
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

func (state *State) AddBullet(x float32, y float32, clientId string, BulletCreateMessage *message.CreateBullet) {
	bulletId := uuid.New().String()
	if player, exists := state.Players.Get(clientId); exists {
		newBullet := object.NewBullet(
			clientId,
			BulletCreateMessage.PlayerId,
			(*player).MagicType,
			x,
			y,
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
	if player, exists := state.Players.Get(clientId); exists {
		state.AddBullet((*(*player).GetPhysical()).GetX(), (*(*player).GetPhysical()).GetY(), clientId, msg)
	}
}

func (state *State) DoDash(clientId string, msg *message.DoDash) {
	if player, exists := state.Players.Get(clientId); exists {
		(*player).DoDash()
	}
}
