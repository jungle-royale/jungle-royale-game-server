package object

import (
	"jungle-royale/message"
	"jungle-royale/physical"
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

func NewMagicItem(magicType int, x, y float64) *Magic {
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
		X:         float32(magic.physicalObject.GetX()),
		Y:         float32(magic.physicalObject.GetY()),
	}
}

func (magic *Magic) GetObjectType() int {
	return OBJECT_MAGICITEM
}

func (magic *Magic) GetObjectId() string {
	return magic.ItemId
}

func (magic *Magic) MakeGetItemState(playerId string) GetItemState {
	if magic.magicType == STONE_MAGIC {
		return NewGetItemState(
			magic.ItemId,
			playerId,
			ITEM_STONE_MAGIC,
			magic.physicalObject.GetX(),
			magic.physicalObject.GetY(),
		)
	} else if magic.magicType == FIRE_MAGIC {
		return NewGetItemState(
			magic.ItemId,
			playerId,
			ITEM_FIRE_MAGIC,
			magic.physicalObject.GetX(),
			magic.physicalObject.GetY(),
		)
	}
	return GetItemState{}
}
