package object

import (
	"jungle-royale/message"
	"jungle-royale/object/physical"
	"math"
	"sync"
	"time"

	"github.com/google/uuid"
)

const PLAYER_SPEED = 0.07

// const PLAYER_SPEED = 1
const DASH_SPEED = 0.3
const DASH_TICK = 6
const DASH_COOLTIME = 10 // 0.1 sec
const PLAYER_RADIOUS = 0.5
const EPSILON = 1e-9
const SHOOTING_COOLTIME = 6 // 0.1 sec

// dying status
const (
	DYING_NONE = iota
	DYING_SNOW
	DYING_STONE
	DYING_FIRE
	DYING_FALL
)

func (pd *PlayerDeadState) Kill() {
	pd.KillNum++
}

type Player struct {
	mu               sync.Mutex
	id               string
	dir              float64 // (dx, dy)
	angle            float64 // degree
	speed            float64
	isMoveing        bool
	isDashing        bool
	dashTime         int
	dashCoolTime     int
	health           int
	MagicType        int
	DyingStatus      *PlayerDeadState
	physicalObject   physical.Physical
	IsShooting       bool
	ShootingCoolTime int
}

func NewPlayer(id string, x, y float64) *Player {

	return &Player{
		sync.Mutex{},
		id,
		0,
		0,
		PLAYER_SPEED,
		false,
		false,
		0,
		0,
		100,
		BULLET_NONE,
		NewPlayerDeadState("", id, DYING_NONE),
		physical.NewCircle(x, y, PLAYER_RADIOUS),
		false,
		0,
	}
}

func (player *Player) CreateBullet() *Bullet {
	if player.ShootingCoolTime > 0 {
		return nil
	} else {
		player.ShootingCoolTime = SHOOTING_COOLTIME
	}
	newBullet := NewBullet(
		uuid.New().String(),
		player.id,
		player.MagicType,
		player.physicalObject.GetX(),
		player.physicalObject.GetY(),
		player.angle,
	)
	return newBullet
}

func (player *Player) SetLocation(x, y float64) {
	player.mu.Lock()
	player.physicalObject.SetCoord(x, y)
	player.mu.Unlock()
}

func (player *Player) CalcGameTick() {
	player.mu.Lock()
	if player.isMoveing {
		player.physicalObject.Move()
	}
	if player.dashCoolTime > 0 {
		player.dashCoolTime--
	}
	if player.dashTime > 0 {
		player.dashTime--
		if player.dashTime == 0 {
			player.isDashing = false
			player.dashCoolTime = DASH_COOLTIME
			player.speed = PLAYER_SPEED
			player.physicalObject.SetDir(
				float64(player.speed)*math.Sin(player.dir*(math.Pi/180)),
				float64(player.speed)*math.Cos(player.dir*(math.Pi/180))*-1,
			)
		}
	}
	if player.ShootingCoolTime > 0 {
		player.ShootingCoolTime--
	}
	player.mu.Unlock()
}

func (player *Player) IsValid() bool {
	return player.health > 0
}

func (player *Player) DirChange(angle float64, isMoved bool) {
	player.mu.Lock()
	player.dir = angle
	player.physicalObject.SetDir(
		player.speed*math.Sin(angle*(math.Pi/180)),
		player.speed*math.Cos(angle*(math.Pi/180))*-1,
	)
	player.isMoveing = isMoved
	player.mu.Unlock()
}

func (player *Player) AngleChange(angle float64) {
	player.angle = angle
}

func (player *Player) MakeSendingData() *message.PlayerState {
	return &message.PlayerState{
		Id:           player.id,
		X:            float32(player.physicalObject.GetX()),
		Y:            float32(player.physicalObject.GetY()),
		Health:       int32(player.health),
		MagicType:    int32(player.MagicType),
		Angle:        float32(player.angle),
		DashCoolTime: int32(player.dashCoolTime),
	}
}

func (player *Player) HitedBullet(bullet *Bullet) bool {

	if bullet.playerId == player.id {
		return false
	}

	if bullet.BulletType == BULLET_NONE {
		player.mu.Lock()
		player.health -= BULLET_DAMAGE
		player.DyingStatus.DyingStatus = DYING_SNOW
		player.DyingStatus.Killer = bullet.playerId
		player.mu.Unlock()
	} else if bullet.BulletType == BULLET_STONE {
		player.mu.Lock()
		player.health -= BULLET_STONE_DAMAGE
		player.DyingStatus.DyingStatus = DYING_STONE
		player.DyingStatus.Killer = bullet.playerId
		player.mu.Unlock()
	} else if bullet.BulletType == BULLET_FIRE {
		player.mu.Lock()
		player.health -= BULLET_DAMAGE
		player.DyingStatus.DyingStatus = DYING_FIRE
		player.DyingStatus.Killer = bullet.playerId
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
			player.DyingStatus.DyingStatus = DYING_FIRE
			player.DyingStatus.Killer = bullet.playerId
			player.mu.Unlock()
			count--
		}
	}

	return true
}

func (player *Player) Dead(killer string, dyingStatus int, placement int) {
	player.DyingStatus.Killer = killer
	player.DyingStatus.DyingStatus = dyingStatus
	player.DyingStatus.Placement = placement
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
	if !player.isDashing && player.dashCoolTime == 0 {
		player.mu.Lock()
		player.isDashing = true
		player.speed = DASH_SPEED
		player.physicalObject.SetDir(
			player.speed*math.Sin(player.dir*(math.Pi/180)),
			player.speed*math.Cos(player.dir*(math.Pi/180))*-1,
		)
		player.dashTime = DASH_TICK
		player.mu.Unlock()
	}
}

func (player *Player) GetPhysical() *physical.Physical {
	return &player.physicalObject
}

func (player *Player) GetObjectType() int {
	return OBJECT_PLAYER
}

func (player *Player) GetObjectId() string {
	return player.id
}
