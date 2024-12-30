package object

import (
	"jungle-royale/message"
	"jungle-royale/object/physical"
	"sync"

	"github.com/google/uuid"
)

// magic type
const (
	NONE_MAGIC = iota
	STONE_MAGIC
	FIRE_MAGIC
)

const MAGIC_ITEM_RADIOUS = 0.3

type Magic struct {
	mu             sync.Mutex
	ItemId         string
	magicType      int
	physicalObject physical.Physical
}

func NewMagicItem(magicType int, x float32, y float32) *Magic {
	return &Magic{
		sync.Mutex{},
		uuid.New().String(),
		magicType,
		physical.NewCircle(x, y, MAGIC_ITEM_RADIOUS),
	}
}

func (magic *Magic) DoEffet(player *Player) {
	player.GetMagic(magic.magicType)
}

func (magic *Magic) GetPhysical() *physical.Physical {
	return &magic.physicalObject
}

func (magic *Magic) MakeSendingData() *message.MagicItemState {
	return &message.MagicItemState{
		ItemId:    magic.ItemId,
		MagicType: int32(magic.magicType),
		X:         magic.physicalObject.GetX(),
		Y:         magic.physicalObject.GetY(),
	}
}
