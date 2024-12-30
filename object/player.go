package object

import (
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object/physical"
	"math"
	"sync"
	"time"
)

const PLAYER_SPEED = 0.3
const DASH_SPEED = 1.2
const DASH_TICK = 5
const DASH_COOLTIME = 60 // 1sec
const PLAYER_RADIOUS = 0.5
const EPSILON = 1e-9

type Player struct {
	mu             sync.Mutex
	id             string
	angle          float64 // degree
	speed          float32
	isMoveing      bool
	isDashing      bool
	dashCoolTime   int
	health         int
	MagicType      int
	collisionList  []int
	physicalObject physical.Physical
}

func NewPlayer(id string, x float32, y float32) *Player {

	return &Player{
		sync.Mutex{},
		id,
		0,
		PLAYER_SPEED,
		false,
		false,
		0,
		100,
		NONE_MAGIC,
		[]int{ObjectHealPack, ObjectMagicItem},
		physical.NewCircle(x, y, PLAYER_RADIOUS),
	}
}

func (player *Player) SetLocation(x float32, y float32) {
	player.mu.Lock()
	player.physicalObject.SetCoord(x, y)
	player.mu.Unlock()
}

func (player *Player) CalcGameTick() {
	player.mu.Lock()
	if player.isMoveing {
		dx := player.speed * float32(math.Sin(player.angle*(math.Pi/180)))
		dy := player.speed * float32(math.Cos(player.angle*(math.Pi/180))) * -1
		player.physicalObject.Move(dx, dy)
		player.dashCoolTime--
	}
	player.mu.Unlock()
}

func (player *Player) IsValid() bool {
	return player.health > 0
}

func (player *Player) DirChange(angle float64, isMoved bool) {
	player.mu.Lock()
	player.angle = angle
	player.isMoveing = isMoved
	player.mu.Unlock()
}

func (player *Player) MakeSendingData() *message.PlayerState {
	return &message.PlayerState{
		Id:        player.id,
		X:         player.physicalObject.GetX(),
		Y:         player.physicalObject.GetY(),
		Health:    int32(player.health),
		MagicType: int32(player.MagicType),
	}
}

func (player *Player) HeatedBullet(bulletType int) {
	if bulletType == BULLET_NONE {
		player.mu.Lock()
		player.health -= BULLET_DAMAGE
		player.mu.Unlock()
	} else if bulletType == BULLET_STONE {
		player.mu.Lock()
		player.health -= BULLET_STONE_DAMAGE
		player.mu.Unlock()
	} else if bulletType == BULLET_FIRE {
		player.mu.Lock()
		player.health -= BULLET_DAMAGE
		player.mu.Unlock()
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		count := BULLET_FIRE_LAST_SEC
		for range ticker.C {
			if count <= 0 {
				break
			}
			player.mu.Lock()
			player.health -= BULLET_FIRE_SEC_DAMAGE
			player.mu.Unlock()
			count--
		}
	}
}

func (player *Player) GetHealPack() {
	player.mu.Lock()
	player.health += HEAL_AMOUNT
	if player.health >= 100 {
		player.health = 100
	}
	player.mu.Unlock()
}

func (player *Player) GetMagic(magicType int) {
	player.mu.Lock()
	player.MagicType = magicType
	player.mu.Unlock()
}

func (player *Player) DoDash() {
	if !player.isDashing && player.dashCoolTime < 0 {
		player.mu.Lock()
		player.isDashing = true
		player.speed = DASH_SPEED
		player.mu.Unlock()
		time.AfterFunc(cons.CalcLoopInterval*DASH_TICK*time.Millisecond, func() {
			player.mu.Lock()
			player.isDashing = false
			player.speed = PLAYER_SPEED
			player.dashCoolTime = DASH_COOLTIME
			player.mu.Unlock()
		})
	}
}

func (player *Player) GetPhysical() *physical.Physical {
	return &player.physicalObject
}
