package object

import (
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object/physical"
)

// tile state
const (
	TILE_NORMAL = iota
	TILE_DANGEROUS
	TILE_FALL
)

type Tile struct {
	TileId         int
	TileState      int
	PhysicalObject *physical.Rectangle
}

func NewTile(tileId int, x, y float32) *Tile {
	return &Tile{
		TileId:    tileId,
		TileState: TILE_NORMAL,
		PhysicalObject: &physical.Rectangle{
			X:      x,
			Y:      y,
			Width:  cons.CHUNK_LENGTH,
			Length: cons.CHUNK_LENGTH,
		},
	}
}

func (tile *Tile) SetTileState(tileState int) {
	tile.TileState = tileState
}

func (tile *Tile) MakeSendingData() *message.TileState {
	return &message.TileState{
		TileId:    int32(tile.TileId),
		TileState: int32(tile.TileState),
		X:         tile.PhysicalObject.X,
		Y:         tile.PhysicalObject.Y,
	}
}
