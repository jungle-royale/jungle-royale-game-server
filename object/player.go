package object

import (
	"jungle-royale/message"
	"jungle-royale/object/physical"
	"math"
	"sync"

	"github.com/google/uuid"
)

const PLAYER_SPEED = 0.06

// const PLAYER_SPEED = 1
const DASH_SPEED = 0.2
const DASH_TICK = 12
const DASH_COOLTIME = 10 // 0.1 sec
const PLAYER_RADIOUS = 0.5
const EPSILON = 1e-9
const SHOOTING_COOLTIME = 10 // 0.1 sec

// dying status
const (
	DYING_NONE = iota
	DYING_SNOW
	DYING_STONE
	DYING_FIRE
	DYING_FALL
)

type Player struct {
	Mu                 sync.Mutex
	id                 string
	dir                float64 // (dx, dy)
	angle              float64 // degree
	speed              float64
	isMoveing          bool
	isDashing          bool
	dashTime           int
	dashCoolTime       int
	health             int
	MagicType          int
	DyingStatus        *PlayerDeadState
	physicalObject     physical.Physical
	IsShooting         bool
	ShootingCoolTime   int
	FireDamageTickTime int
	FireDamageCount    int
	FireDamageOwner    string // 불 틱데미지 입힌사람
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
		0,
		0,
		"",
	}
}

func (player *Player) CreateBullet() *Bullet {
	if player.ShootingCoolTime > 0 || player.isDashing {
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

func (player *Player) Kill() {
	player.Mu.Lock()
	player.DyingStatus.KillNum++
	player.Mu.Unlock()
}

func (player *Player) SetLocation(x, y float64) {
	player.Mu.Lock()
	player.physicalObject.SetCoord(x, y)
	player.Mu.Unlock()
}

func (player *Player) CalcGameTick() {
	player.Mu.Lock()
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
	if player.FireDamageCount > 0 {
		player.FireDamageTickTime--
		if player.FireDamageTickTime == 0 {
			player.health -= BULLET_FIRE_SEC_DAMAGE
			player.FireDamageTickTime = BULLET_FIRE_LAST_TICK
			player.FireDamageCount--
			player.DyingStatus.Killer = player.FireDamageOwner
		}
	}
	player.Mu.Unlock()
}

func (player *Player) IsValid() bool {
	return player.health > 0
}

func (player *Player) DirChange(angle float64, isMoved bool) {
	player.Mu.Lock()
	player.dir = angle
	player.physicalObject.SetDir(
		player.speed*math.Sin(angle*(math.Pi/180)),
		player.speed*math.Cos(angle*(math.Pi/180))*-1,
	)
	player.isMoveing = isMoved
	player.Mu.Unlock()
}

func (player *Player) AngleChange(angle float64) {
	player.angle = angle
}

func (player *Player) MakeSendingData() *message.PlayerState {
	burn := false
	if player.FireDamageCount > 0 {
		burn = true
	}
	return &message.PlayerState{
		Id:           player.id,
		X:            float32(player.physicalObject.GetX()),
		Y:            float32(player.physicalObject.GetY()),
		Health:       int32(player.health),
		MagicType:    int32(player.MagicType),
		Angle:        float32(player.angle),
		DashCoolTime: int32(player.dashCoolTime),
		IsMoved:      player.isMoveing,
		IsDashing:    player.isDashing,
		IsBurn:       burn,
	}
}

func (player *Player) HitedBullet(bullet *Bullet) bool {

	if bullet.playerId == player.id {
		return false
	}
	if bullet.BulletType == BULLET_NONE {
		player.Mu.Lock()
		player.health -= BULLET_DAMAGE
		player.DyingStatus.DyingStatus = DYING_SNOW
		player.DyingStatus.Killer = bullet.playerId
		player.Mu.Unlock()
	} else if bullet.BulletType == BULLET_STONE {
		player.Mu.Lock()
		player.health -= BULLET_STONE_DAMAGE
		player.DyingStatus.DyingStatus = DYING_STONE
		player.DyingStatus.Killer = bullet.playerId
		player.Mu.Unlock()
	} else if bullet.BulletType == BULLET_FIRE {
		player.Mu.Lock()
		player.health -= BULLET_DAMAGE
		player.DyingStatus.DyingStatus = DYING_FIRE
		player.DyingStatus.Killer = bullet.playerId
		player.FireDamageTickTime = BULLET_FIRE_LAST_TICK
		player.FireDamageCount = BULLET_FIRE_LAST_COUNT
		player.FireDamageOwner = bullet.playerId
		player.Mu.Unlock()
	}

	return true
}

func (player *Player) Dead(killer string, dyingStatus int, placement int) {
	player.Mu.Lock()
	player.DyingStatus.Killer = killer
	player.DyingStatus.DyingStatus = dyingStatus
	player.DyingStatus.Placement = placement
	player.Mu.Unlock()
}

func (player *Player) GetHealPack() {
	player.Mu.Lock()
	player.health += HEAL_AMOUNT
	if player.health >= 100 {
		player.health = 100
	}
	player.Mu.Unlock()
}

func (player *Player) GetMagic(magicType int) {
	player.Mu.Lock()
	player.MagicType = magicType
	player.Mu.Unlock()
}

func (player *Player) DoDash() bool {
	if !player.isDashing && player.dashCoolTime == 0 {
		player.Mu.Lock()
		player.isDashing = true
		player.speed = DASH_SPEED
		player.physicalObject.SetDir(
			player.speed*math.Sin(player.dir*(math.Pi/180)),
			player.speed*math.Cos(player.dir*(math.Pi/180))*-1,
		)
		player.dashTime = DASH_TICK
		player.Mu.Unlock()
		return true
	} else {
		return false
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
