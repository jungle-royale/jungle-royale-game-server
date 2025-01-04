package object

import (
	"jungle-royale/message"
	"jungle-royale/object/physical"
	"math"
	"sync"
)

const BULLET_SPEED = 0.2
const BULLET_RANGE = 20.0
const BULLET_DAMAGE = 20
const BULLET_RADIOUS = 0.2
const BULLET_MAX_TICK = BULLET_RANGE / BULLET_SPEED

// bullet type (= magic type)
const (
	BULLET_NONE = iota
	BULLET_STONE
	BULLET_FIRE
)

const BULLET_STONE_DAMAGE = 30
const BULLET_FIRE_SEC_DAMAGE = 3
const BULLET_FIRE_LAST_SEC = 5

type Bullet struct {
	mu             sync.Mutex
	bulletId       string
	playerId       string
	dx             float32
	dy             float32
	lastTick       int
	BulletType     int
	isValid        bool
	physicalObject physical.Physical
}

func NewBullet(
	bulletId string,
	playerId string,
	magicType int,
	startX float32,
	startY float32,
	angle float64,
) *Bullet {
	dx := float32(BULLET_SPEED * math.Sin(angle*(math.Pi/180)))
	dy := -1 * float32(BULLET_SPEED*math.Cos(angle*(math.Pi/180)))
	return &Bullet{
		sync.Mutex{},
		bulletId,
		playerId,
		dx,
		dy,
		BULLET_MAX_TICK,
		magicType,
		true,
		physical.NewCircle(startX+dx, startY+dy, BULLET_RADIOUS),
	}
}

func (bullet *Bullet) CalcGameTick() {
	bullet.mu.Lock()
	bullet.physicalObject.Move(bullet.dx, bullet.dy)
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
		BulletId: bullet.bulletId,
		X:        bullet.physicalObject.GetX(),
		Y:        bullet.physicalObject.GetY(),
	}
}

func (bullet *Bullet) GetPhysical() *physical.Physical {
	return &bullet.physicalObject
}

func (bullet *Bullet) GetObjectType() int {
	return OBJECT_BULLET
}

func (bullet *Bullet) GetObjectId() string {
	return bullet.bulletId
}

func (bullet *Bullet) MakeHeatBulletState(heatedPlayerId string) HeatBulletState {
	return NewHeatBulletState(bullet.bulletId, heatedPlayerId)
}
