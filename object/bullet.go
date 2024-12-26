package object

import (
	"jungle-royale/message"
	"log"
	"math"
)

const BULLET_SPEED = 20.0
const BULLET_RANGE = 500.0
const BULLET_DAMAGE = 20
const BULLET_MAX_TICK = BULLET_RANGE / BULLET_SPEED

type Bullet struct {
	bulletId      string
	playerId      string
	x             float32
	y             float32
	dx            float32
	dy            float32
	lastTick      int
	collisionList map[string]*Player
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
		bulletId,
		playerId,
		startX,
		startY,
		dx,
		dy,
		BULLET_MAX_TICK,
		collisionList,
	}
}

func (bullet *Bullet) Collision() bool {
	for id := range bullet.collisionList {
		player := bullet.collisionList[id]
		if bullet.playerId == player.id {
			continue
		}

		distanceSquared := math.Pow(float64(bullet.x-player.x), 2) + math.Pow(float64(bullet.y-player.y), 2)
		radiusSquared := math.Pow(float64(player.radious), 2)
		if radiusSquared > distanceSquared {
			log.Printf("collision with %s %s", bullet.playerId, player.id)
			player.HeatedBullet()
			return true
		}
	}
	return false
}

func (bullet *Bullet) Move() bool {
	bullet.x += bullet.dx
	bullet.y += bullet.dy
	bullet.lastTick--

	if bullet.Collision() {
		return false
	}

	if bullet.lastTick > 0 {
		return true
	} else {
		return false
	}
}

func (bullet *Bullet) MakeSendingData() *message.BulletState {
	return &message.BulletState{
		BulletId: bullet.bulletId,
		X:        bullet.x,
		Y:        bullet.y,
	}
}
