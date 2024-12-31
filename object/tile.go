package object

import (
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object/physical"
)

type Tile struct {
	TileId         int
	PhysicalObject *physical.Rectangle
}

func NewTile(tileId int, x, y float32) *Tile {
	return &Tile{
		TileId: tileId,
		PhysicalObject: &physical.Rectangle{
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
		X:      tile.PhysicalObject.X,
		Y:      tile.PhysicalObject.Y,
	}
}
