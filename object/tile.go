package object

import (
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object/physical"
)

type Tile struct {
	TileId         int
	physicalObject physical.Rectangle
}

func NewTile(tileId int, x, y float32) *Tile {
	return &Tile{
		TileId: tileId,
		physicalObject: physical.Rectangle{
			X:      x,
			Y:      y,
			Width:  cons.CHUNK_LENGTH,
			Length: cons.CHUNK_LENGTH,
		},
	}
}

func (tile *Tile) MakeSendingData() *message.TileState {
	return &message.TileState{
		TileId: int32(tile.TileId),
		X:      tile.physicalObject.X,
		Y:      tile.physicalObject.Y,
	}
}
