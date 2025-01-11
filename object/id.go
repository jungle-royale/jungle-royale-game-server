package object

import "sync"

type ObjectIdAllocator struct {
	playerId        int
	playerIdMu      sync.Mutex
	bulletId        int
	bulletIdMu      sync.Mutex
	healPackId      int
	healPackIdMu    sync.Mutex
	magicId         int
	magicIdMu       sync.Mutex
	tileId          int
	tileIdMu        sync.Mutex
	environmentId   int
	environmentIdMu sync.Mutex
}

func NewObjectIdAllocator() *ObjectIdAllocator {
	return &ObjectIdAllocator{
		0,
		sync.Mutex{},
		0,
		sync.Mutex{},
		0,
		sync.Mutex{},
		0,
		sync.Mutex{},
		0,
		sync.Mutex{},
		0,
		sync.Mutex{},
	}
}

func (oia *ObjectIdAllocator) AllocatePlayerId() int {
	oia.playerIdMu.Lock()
	id := oia.playerId
	oia.playerId++
	oia.playerIdMu.Unlock()
	return id
}

func (oia *ObjectIdAllocator) AllocateBulletId() int {
	oia.bulletIdMu.Lock()
	id := oia.bulletId
	oia.bulletId++
	oia.bulletIdMu.Unlock()
	return id
}

func (oia *ObjectIdAllocator) AllocateHealPackId() int {
	oia.healPackIdMu.Lock()
	id := oia.healPackId
	oia.healPackId++
	oia.healPackIdMu.Unlock()
	return id
}

func (oia *ObjectIdAllocator) AllocateMagicId() int {
	oia.magicIdMu.Lock()
	id := oia.magicId
	oia.magicId++
	oia.magicIdMu.Unlock()
	return id
}

func (oia *ObjectIdAllocator) AllocateTileId() int {
	oia.tileIdMu.Lock()
	id := oia.tileId
	oia.tileId++
	oia.tileIdMu.Unlock()
	return id
}

func (oia *ObjectIdAllocator) AllocateEnvironmentId() int {
	oia.environmentIdMu.Lock()
	id := oia.environmentId
	oia.environmentId++
	oia.environmentIdMu.Unlock()
	return id
}
