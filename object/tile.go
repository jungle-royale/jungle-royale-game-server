package object

import (
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object/physical"
	"jungle-royale/util"
)

// tile state
const (
	TILE_NORMAL = iota
	TILE_DANGEROUS
	TILE_FALL
)

type Tile struct {
	IdxI           int
	IdxJ           int
	TileId         string
	TileState      int
	PhysicalObject *physical.Rectangle
	ChildTile      *util.Set[*Tile]
	ParentTile     *Tile
}

func NewTile(tileId string, x, y float32, idxi, idxj int) *Tile {
	return &Tile{
		IdxI:      idxi,
		IdxJ:      idxj,
		TileId:    tileId,
		TileState: TILE_NORMAL,
		PhysicalObject: &physical.Rectangle{
			X:      x,
			Y:      y,
			Width:  cons.CHUNK_LENGTH,
			Length: cons.CHUNK_LENGTH,
		},
		ChildTile:  util.NewSet[*Tile](),
		ParentTile: nil,
	}
}

func (tile *Tile) SetTileState(tileState int) {
	tile.TileState = tileState
}

func (tile *Tile) MakeSendingData() *message.TileState {
	return &message.TileState{
		TileId:    tile.TileId,
		TileState: int32(tile.TileState),
		X:         tile.PhysicalObject.X,
		Y:         tile.PhysicalObject.Y,
	}
}
