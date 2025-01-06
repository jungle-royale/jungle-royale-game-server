package object

import (
	"jungle-royale/object/physical"

	"github.com/google/uuid"
)

const TILE_TYPE_NUM = 0

type objectData interface {
	createEnvObject(dx, dy float64) *EnvObject
}

type objectCircle struct {
	x       float64
	y       float64
	radious float64
	isShort bool
}

func (c objectCircle) createEnvObject(dx, dy float64) *EnvObject {
	var p physical.Physical = physical.NewCircle(c.x+dx, c.y+dy, c.radious)
	return &EnvObject{
		objId:          uuid.New().String(),
		physicalObject: p,
	}
}

type objectRectangle struct {
	x       float64
	y       float64
	width   float64
	length  float64
	isShort bool
}

func (r objectRectangle) createEnvObject(dx, dy float64) *EnvObject {
	var p physical.Physical = physical.NewRectangle(r.x+dx, r.y+dy, r.width, r.length)
	return &EnvObject{
		objId:          uuid.New().String(),
		physicalObject: p,
	}
}

type EnvObject struct {
	objId          string
	IsShort        bool // false → collision with bullet
	physicalObject physical.Physical
}

func (eo *EnvObject) GetObjectType() int {
	return OBJECT_ENVIRONMENT
}

func (eo *EnvObject) GetObjectId() string {
	return eo.objId
}

func (eo *EnvObject) GetPhysical() *physical.Physical {
	return &eo.physicalObject
}

// [tile type][number of objects]
var environment = [][]objectData{
	// tile type 0
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
		objectRectangle{0.5, 9, 3, 2, false},
		objectCircle{5, 18, 0.5, false},
		objectCircle{8, 7, 0.7, true},
	},

	// tile type 1
	{
		objectRectangle{10.5, 17.0, 3.0, 2.0, false},
		objectCircle{15.0, 13.0, 1.0, false},
		objectCircle{4.0, 15.0, 1.0, false},
		objectCircle{3.0, 5.0, 0.8, false},
		objectCircle{17.0, 4.0, 1.0, false},
		objectCircle{10.0, 8.0, 0.6, false},
		objectRectangle{6.5, 7.94, 1.7, 0.8, true},
	},
}

func (tile *Tile) SetTileEnvironment(tileType int, dx, dy float64) {
	tile.tileType = tileType
	for _, p := range environment[tileType] {
		tile.Environment.Add(p.createEnvObject(dx, dy))
	}
}
