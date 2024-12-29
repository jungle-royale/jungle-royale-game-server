package object

import (
	"jungle-royale/object/physical"
	"reflect"
	"sync"

	"google.golang.org/protobuf/proto"
)

// Mover object enum
const OBJECT_NUM = 2
const MOVER_OBJECT_NUM = 2
const (
	// mover object
	ObjectPlayer = iota
	ObjectBullet

	// nonmover object
	ObjectHealPack
	ObjectMagicItem
)

type Collider interface {
	getPhysical() *physical.Physical
}

type Mover interface {
	CalcGameTick() // move, collision
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

type SyncMapList struct {
	objectLists map[int]*ObjectSyncMap
}

func NewMoverSyncMapList() *SyncMapList {

	list := make(map[int]*ObjectSyncMap)
	list[ObjectPlayer] = NewObjectSyncMap(reflect.TypeOf(Player{}))
	list[ObjectBullet] = NewObjectSyncMap(reflect.TypeOf(Bullet{}))
	list[ObjectHealPack] = NewObjectSyncMap((reflect.TypeOf(HealPack{})))
	list[ObjectMagicItem] = NewObjectSyncMap((reflect.TypeOf(Magic{})))

	return &SyncMapList{list}
}

func (mlist *SyncMapList) GetPlayers() *sync.Map {
	return &mlist.objectLists[ObjectPlayer].Map
}

func (mlist *SyncMapList) GetBullets() *sync.Map {
	return &mlist.objectLists[ObjectBullet].Map
}

func (mlist *SyncMapList) GetHealPacks() *sync.Map {
	return &mlist.objectLists[ObjectHealPack].Map
}

func (mlist *SyncMapList) GetMagicItems() *sync.Map {
	return &mlist.objectLists[ObjectMagicItem].Map
}

type Item interface {
	DoEffet(p *Player)
}

type Wepone interface {
	GiveItem(p *Player)
}
