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
	collisionList  []int
	physicalObject physical.Physical
}

func NewBullet(
	bulletId string,
	playerId string,
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
		true,
		[]int{ObjectPlayer},
		physical.NewCircle(startX+dx, startY+dy, BULLET_RADIOUS),
	}
}

func (bullet *Bullet) CalcCollision(objectMapList *SyncMapList) *Collider {
	var ret Collider
	flag := false
	for _, c := range bullet.collisionList {
		objectMapList.objectLists[c].Map.Range(func(key, value any) bool {
			switch v := value.(type) {
			case *Player:
				if bullet.physicalObject.IsCollide((*v).getPhysical()) {
					ret = v
					flag = true
					return false
				} else {
					return true
				}
			default:
				return true
			}
		})
	}

	if flag {
		return &ret
	} else {
		return nil
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

func (bullet *Bullet) getPhysical() *physical.Physical {
	return &bullet.physicalObject
}
