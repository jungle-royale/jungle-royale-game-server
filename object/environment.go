package object

import (
	"jungle-royale/physical"
)

type objectData interface {
	createEnvObject(dx, dy float64, id int) *EnvObject
}

type objectCircle struct {
	x       float64
	y       float64
	radious float64
	isShort bool
}

func (c objectCircle) createEnvObject(dx, dy float64, id int) *EnvObject {
	var p physical.Physical = physical.NewCircle(c.x+dx, c.y+dy, c.radious)
	return &EnvObject{
		objId:          id,
		physicalObject: p,
		IsShort:        c.isShort,
	}
}

type objectRectangle struct {
	x       float64
	y       float64
	width   float64
	length  float64
	isShort bool
}

func (r objectRectangle) createEnvObject(dx, dy float64, id int) *EnvObject {
	var p physical.Physical = physical.NewRectangle(r.x+dx, r.y+dy, r.width, r.length)
	return &EnvObject{
		objId:          id,
		physicalObject: p,
		IsShort:        r.isShort,
	}
}

type EnvObject struct {
	objId          int
	IsShort        bool // false → collision with bullet
	physicalObject physical.Physical
}

func (eo *EnvObject) GetObjectType() int {
	return OBJECT_ENVIRONMENT
}

func (eo *EnvObject) GetObjectId() int {
	return eo.objId
}

func (eo *EnvObject) GetPhysical() *physical.Physical {
	return &eo.physicalObject
}

const TILE_TYPE_NUM = 4

// [tile type][number of objects]
var environment = [][]objectData{

	// tile type 0
	{
		objectCircle{18, 18, 1, false},
		objectCircle{2, 18, 1, false},
		objectCircle{2, 2, 1, false},
		objectCircle{18, 2, 1, false},
		objectCircle{4, 1, 0.8, false},
		objectCircle{3, 15, 1, false},
		objectCircle{15, 19, 0.8, false},
		objectRectangle{13.125, 0.6000004, 1.75, 0.8, false},
		objectCircle{6, 5, 0.6, false},
		objectRectangle{7.5, 17, 3, 2, false},
		objectRectangle{9.5, 7.75, 1, 0.5, false},
		objectRectangle{10.5, 7.75, 1, 0.5, false},
		objectCircle{10, 10, 1.4, true},
		objectCircle{17, 5, 0.2, false},
		objectCircle{4, 5, 0.2, false},
		objectCircle{17, 15, 0.2, false},
		objectCircle{6, 16, 0.2, false},
	},

	// tile type 1
	{
		objectCircle{11, 10, 0.2, false},
		objectCircle{11, 9, 0.5, false},
		objectCircle{10, 10, 0.7, false},
		objectCircle{6, 14, 0.8, false},
		objectCircle{14, 14, 0.8, false},
		objectCircle{6, 6, 0.8, false},
		objectCircle{14, 6, 0.8, false},
		objectCircle{9, 1, 0.6, false},
	},

	// tile type 2
	{
		objectCircle{10, 10, 0.6, false},
		objectCircle{15, 15, 0.6, false},
		objectCircle{5, 15, 0.6, false},
		objectCircle{5, 5, 0.6, false},
		objectCircle{15, 5, 0.6, false},
	},

	// tile type 3
	{
		objectCircle{16, 16, 1, false},
		objectCircle{4, 4, 1, false},
		objectCircle{16, 4, 1, false},
		objectCircle{4, 16, 1, false},
		objectRectangle{8.5, 9.25, 3, 1.5, false},
	},

	// empty tile (type 4)
	{},
}

func (tile *Tile) SetTileEnvironment(tileType, objectId int, dx, dy float64) {
	tile.tileType = tileType
	for _, p := range environment[tileType] {
		tile.Environment.Add(p.createEnvObject(dx, dy, objectId))
	}
}
