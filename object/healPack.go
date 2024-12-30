package object

import (
	"jungle-royale/message"
	"jungle-royale/object/physical"
	"sync"

	"github.com/google/uuid"
)

const HEAL_AMOUNT = 30
const HEALPACK_WIDTH = 0.5
const HEALPACK_LENGTH = 0.3

type HealPack struct {
	mu             sync.Mutex
	Id             string
	physicalObject physical.Physical
}

func NewHealPack(x float32, y float32) *HealPack {
	return &HealPack{
		sync.Mutex{},
		uuid.New().String(),
		physical.NewRectangle(x, y, HEALPACK_WIDTH, HEALPACK_LENGTH),
	}
}

func (heal *HealPack) DoEffet(player *Player) {
	player.GetHealPack()
}

func (heal *HealPack) GetPhysical() *physical.Physical {
	return &heal.physicalObject
}

func (heal *HealPack) SetLocation(x float32, y float32) {
	heal.mu.Lock()
	heal.physicalObject.SetCoord(x, y)
	heal.mu.Unlock()
}

func (heal *HealPack) MakeSendingData() *message.HealPackState {
	return &message.HealPackState{
		ItemId: heal.Id,
		X:      heal.physicalObject.GetX(),
		Y:      heal.physicalObject.GetY(),
	}
}
