package object

import (
	"jungle-royale/message"
	"jungle-royale/physical"
	"math"
	"sync"
)

const BULLET_SPEED = 0.25
const BULLET_DAMAGE = 5
const BULLET_RADIOUS = 0.2
const BULLET_MAX_TICK = 60

// bullet type (= magic type)
const (
	BULLET_NONE = iota
	BULLET_STONE
	BULLET_FIRE
)

const BULLET_STONE_DAMAGE = BULLET_DAMAGE + 3
const BULLET_FIRE_SEC_DAMAGE = 1
const BULLET_FIRE_LAST_TICK = 60
const BULLET_FIRE_LAST_COUNT = 5

type Bullet struct {
	mu             sync.Mutex
	bulletId       int
	playerId       int
	lastTick       int
	BulletType     int
	isValid        bool
	physicalObject physical.Physical
}

func NewBullet(
	bulletId int,
	playerId int,
	magicType int,
	startX float64,
	startY float64,
	angle float64,
) *Bullet {
	dx := BULLET_SPEED * math.Sin(angle*(math.Pi/180))
	dy := -1 * BULLET_SPEED * math.Cos(angle*(math.Pi/180))
	p := physical.NewCircle(startX+dx, startY+dy, BULLET_RADIOUS)
	p.SetDir(dx, dy)
	return &Bullet{
		sync.Mutex{},
		bulletId,
		playerId,
		BULLET_MAX_TICK,
		magicType,
		true,
		p,
	}
}

func (bullet *Bullet) CalcGameTick() {
	bullet.mu.Lock()
	bullet.physicalObject.Move()
	bullet.lastTick--
	if bullet.lastTick <= 0 {
		bullet.isValid = false
	}
	bullet.mu.Unlock()
}

func (bullet *Bullet) IsValid() bool {
	return bullet.isValid
}

func (bullet *Bullet) MakeSendingData() *message.BulletState {
	return &message.BulletState{
		BulletId: int32(bullet.bulletId),
		X:        float32(bullet.physicalObject.GetX()),
		Y:        float32(bullet.physicalObject.GetY()),
	}
}

func (bullet *Bullet) GetPhysical() *physical.Physical {
	return &bullet.physicalObject
}

func (bullet *Bullet) GetObjectType() int {
	return OBJECT_BULLET
}

func (bullet *Bullet) GetObjectId() int {
	return bullet.bulletId
}

func (bullet *Bullet) MakeHitBulletState(objectType int, objectId int) HitBulletState {
	return NewHitBulletState(
		bullet.bulletId,
		objectId,
		objectType,
		bullet.physicalObject.GetX(),
		bullet.physicalObject.GetY(),
		bullet.BulletType,
	)
}

func (bullet *Bullet) IsValidHit(playerId int) bool {
	if bullet.playerId != playerId {
		return true
	} else {
		return false
	}
}
