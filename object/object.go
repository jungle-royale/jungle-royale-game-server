package object

import (
	"reflect"
	"sync"

	"google.golang.org/protobuf/proto"
)

// Mover object enum
const MOVER_NUM = 2
const (
	ObjectPlayer = iota
	ObjectBullet
)

// Non mover object enum
const NONMOVER_NUM = 0
const ()

type Mover interface {
	CalcGameTick() // move, collision
	CalcCollision()
	MakeSendingData() *proto.Message
	IsValid() bool
}

type NonMover interface {
}

type ObjectSyncMap struct {
	ObjectType reflect.Type
	Map        sync.Map
}

// ex) NewMoverSyncMap(Player{})
func NewObjectSyncMap(t reflect.Type) *ObjectSyncMap {
	return &ObjectSyncMap{
		t,
		sync.Map{},
	}
}

type MoverSyncMapList struct {
	list []*ObjectSyncMap
}

func NewMoverSyncMapList() *MoverSyncMapList {

	list := make([]*ObjectSyncMap, 0)
	list = append(list, NewObjectSyncMap(reflect.TypeOf(Player{})))
	list = append(list, NewObjectSyncMap(reflect.TypeOf(Bullet{})))

	return &MoverSyncMapList{list}
}

func (mlist *MoverSyncMapList) GetPlayers() *sync.Map {
	return &mlist.list[ObjectPlayer].Map
}

func (mlist *MoverSyncMapList) GetBullets() *sync.Map {
	return &mlist.list[ObjectBullet].Map
}
