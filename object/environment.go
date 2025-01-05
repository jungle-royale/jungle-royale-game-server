package object

import "jungle-royale/object/physical"

const TILE_TYPE_NUM = 0

type objectData interface {
	createPhysical(dx, dy float64) *physical.Physical
}

type objectCircle struct {
	objectType int // 0: circle, 1: rectangle
	x          float64
	y          float64
	radious    float64
	isShort    bool
}

func (c objectCircle) createPhysical(dx, dy float64) *physical.Physical {
	var ret physical.Physical = physical.NewCircle(c.x+dx, c.y+dy, c.radious)
	return &ret
}

type objectRectangle struct {
	objectType int // 0: circle, 1: rectangle
	x          float64
	y          float64
	width      float64
	length     float64
	isShort    bool
}

func (r objectRectangle) createPhysical(dx, dy float64) *physical.Physical {
	var ret physical.Physical = physical.NewRectangle(r.x+dx, r.y+dy, r.width, r.length)
	return &ret
}

// [tile type][number of objects]
var environment = [][]objectData{
	{objectCircle{0, 1, 1, 1, true}, objectCircle{0, 4, 3, 2, true}},
}

func (tile *Tile) SetTileEnvironment(tileType int, dx, dy float64) {
	tile.tileType = tileType
	for _, p := range environment[tileType] {
		tile.Environment.Add(p.createPhysical(dx, dy))
	}
}
