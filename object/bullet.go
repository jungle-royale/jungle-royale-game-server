package object

import (
	"jungle-royale/message"
	"jungle-royale/object/physical"
	"math"
	"sync"
)

const BULLET_SPEED = 1
const BULLET_RANGE = 10.0
const BULLET_DAMAGE = 20
const BULLET_RADIOUS = 0.2
const BULLET_MAX_TICK = BULLET_RANGE / BULLET_SPEED

type Bullet struct {
	mu             sync.Mutex
	bulletId       string
	playerId       string
	dx             float32
	dy             float32
	lastTick       int
	isValid        bool
	collisionList  map[string]*Player
	physicalObject physical.Physical
}

func NewBullet(
	bulletId string,
	playerId string,
	startX float32,
	startY float32,
	angle float64,
	collisionList map[string]*Player,
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
		true,
		collisionList,
		physical.NewCircle(startX, startY, BULLET_RADIOUS),
	}
}

func (bullet *Bullet) CalcCollision() bool {
	for id := range bullet.collisionList {
		player := bullet.collisionList[id]
		if bullet.playerId == player.id {
			continue
		}

		if bullet.physicalObject.Collide(player.physicalObject) {
			player.health -= BULLET_DAMAGE
			return true
		}
	}
	return false
}

func (bullet *Bullet) CalcGameTick() {
	bullet.mu.Lock()
	bullet.physicalObject.Move(bullet.dx, bullet.dy)
	bullet.lastTick--
	if bullet.CalcCollision() {
		bullet.isValid = false
	}
	if bullet.lastTick < 0 {
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
