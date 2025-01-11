package object

import (
	"jungle-royale/physical"
)

// Mover object enum
const OBJECT_NUM = 4
const MOVER_OBJECT_NUM = 2
const (
	// mover object
	OBJECT_PLAYER = iota
	OBJECT_BULLET

	// nonmover object
	OBJECT_HEALPACK
	OBJECT_MAGICITEM
	OBJECT_ENVIRONMENT
)

type Object interface {
	GetObjectType() int
	GetObjectId() int
	Collider
}

type Collider interface {
	GetPhysical() *physical.Physical
}

type Mover interface {
	CalcGameTick() // move, collision
	IsValid() bool
	Object
}

type NonMover interface {
}
