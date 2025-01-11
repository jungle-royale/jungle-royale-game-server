package object

import (
	"jungle-royale/message"
	"jungle-royale/physical"
	"sync"
)

const HEAL_AMOUNT = 30
const HEALPACK_WIDTH = 0.5
const HEALPACK_LENGTH = 0.3

type HealPack struct {
	mu             sync.Mutex
	Id             int
	physicalObject physical.Physical
}

func NewHealPack(x, y float64, id int) *HealPack {
	return &HealPack{
		sync.Mutex{},
		id,
		physical.NewRectangle(x, y, HEALPACK_WIDTH, HEALPACK_LENGTH),
	}
}

func (heal *HealPack) DoEffet(player *Player) {
	player.GetHealPack()
}

func (heal *HealPack) GetPhysical() *physical.Physical {
	return &heal.physicalObject
}

func (heal *HealPack) SetLocation(x, y float64) {
	heal.mu.Lock()
	heal.physicalObject.SetCoord(x, y)
	heal.mu.Unlock()
}

func (heal *HealPack) MakeSendingData() *message.HealPackState {
	return &message.HealPackState{
		ItemId: int32(heal.Id),
		X:      float32(heal.physicalObject.GetX()),
		Y:      float32(heal.physicalObject.GetY()),
	}
}

func (heal *HealPack) GetObjectType() int {
	return OBJECT_HEALPACK
}

func (heal *HealPack) GetObjectId() int {
	return heal.Id
}
