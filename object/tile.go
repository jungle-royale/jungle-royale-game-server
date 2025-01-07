package object

import (
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object/physical"
	"jungle-royale/util"
	"sync"
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
	tileType       int
	Environment    *util.Set[*EnvObject]
	Depth          int
	Mu             sync.Mutex
}

func NewTile(tileId string, x, y float64, idxi, idxj int) *Tile {
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
		ChildTile:   util.NewSyncSet[*Tile](),
		ParentTile:  nil,
		Environment: util.NewSyncSet[*EnvObject](),
		Mu:          sync.Mutex{},
	}
}

func (tile *Tile) SetTileState(tileState int) *Tile {
	tile.TileState = tileState
	return tile
}

func (tile *Tile) SetTileType(tileType int) *Tile {
	tile.SetTileEnvironment(
		tileType,
		float64(tile.IdxI*cons.CHUNK_LENGTH),
		float64(tile.IdxJ*cons.CHUNK_LENGTH),
	)
	return tile
}

func (tile *Tile) MakeSendingData() *message.TileState {
	return &message.TileState{
		TileId:    tile.TileId,
		TileState: int32(tile.TileState),
		TileType:  int32(tile.tileType),
		X:         float32(tile.PhysicalObject.X),
		Y:         float32(tile.PhysicalObject.Y),
	}
}

type TileHeap []*Tile

func (th TileHeap) Len() int {
	return len(th)
}

func (th TileHeap) Less(i, j int) bool {
	return th[i].Depth > th[j].Depth
}

func (th TileHeap) Swap(i, j int) {
	th[i], th[j] = th[j], th[i]
}

func (th *TileHeap) Push(e interface{}) {
	*th = append(*th, e.(*Tile))
}

func (th *TileHeap) Pop() interface{} {
	old := *th
	n := len(old)
	e := old[n-1]
	*th = old[0 : n-1]
	return e
}
