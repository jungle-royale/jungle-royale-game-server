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
		objectRectangle{9.5, 7.75, 1.0, 0.5, false},
		objectRectangle{10.5, 7.75, 1.0, 0.5, false},
		objectCircle{10.0, 10.0, 1.4, true},
		objectCircle{3, 4, 0.6, false},
		objectCircle{13.0, 18.0, 1.5, false},
		objectCircle{16.0, 8, 1.3, false},
		objectRectangle{15.75, 17.5, 2.5, 1.0, false},
		objectCircle{8.0, 17.0, 1.7, false},
		objectRectangle{1.15, 17.5, 2.5, 1.0, false},
		objectCircle{15.0, 9.0, 0.2, false},
		objectCircle{11.0, 19.0, 0.2, false},
		objectCircle{18.0, 9.0, 0.2, false},
		objectCircle{2.0, 12.0, 0.2, false},
		objectCircle{2.0, 7.0, 0.2, false},
		objectCircle{18.0, 18.0, 0.2, false},
		objectCircle{2.0, 18.0, 0.2, false},
		objectCircle{8.0, 2.0, 0.2, false},
	},

	// tile type 1
	{
		objectCircle{2, 18, 1, false},
		objectCircle{17, 18, 1, false},
		objectCircle{4, 4, 1, false},
		objectCircle{14, 4, 0.9, false},
		objectCircle{9, 16, 0.8, false},
		objectCircle{12, 12, 1, false},
		objectCircle{18, 11, 0.9, false},
		objectCircle{2, 1, 0.8, false},
		objectRectangle{17.125, 0.6, 1.75, 0.8, false},
		objectCircle{7, 2, 0.6, false},
		objectRectangle{1, 9, 2, 2, false},
		objectRectangle{3, 9, 1, 2, true},
		objectCircle{5, 18, 0.5, false},
		objectCircle{8, 7, 0.7, true},
	},

	// tile type 2
	{
		objectRectangle{10.5, 17.0, 3.0, 2.0, false},
		objectCircle{15.0, 13.0, 1.0, false},
		objectCircle{4.0, 15.0, 1.0, false},
		objectCircle{3.0, 5.0, 0.8, false},
		objectCircle{17.0, 4.0, 1.0, false},
		objectCircle{10.0, 8.0, 0.6, false},
		objectRectangle{16.15, 8.6, 1.7, 0.8, true},
	},

	// tile type 3
	{},
}

func (tile *Tile) SetTileEnvironment(tileType, objectId int, dx, dy float64) {
	tile.tileType = tileType
	for _, p := range environment[tileType] {
		tile.Environment.Add(p.createEnvObject(dx, dy, objectId))
	}
}
