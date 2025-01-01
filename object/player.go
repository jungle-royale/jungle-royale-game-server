package object

import (
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

// dying status
const (
	DYING_NONE = iota
	DYING_SNOW
	DYING_STONE
	DYING_FIRE
	DYING_FALL
)

type PlayerDead struct {
	Killer      string
	dead        string
	DyingStatus int
	KillNum     int
	Placement   int
}

func NewPlayerDead(killer string, dead string, ds int) *PlayerDead {
	return &PlayerDead{
		killer,
		dead,
		ds,
		0,
		-1,
	}
}

func (pd *PlayerDead) Kill() {
	pd.KillNum++
}

func (pd *PlayerDead) MakeSendingData() *message.PlayerDeadState {
	return &message.PlayerDeadState{
		KillerId:    pd.Killer,
		DeadId:      pd.dead,
		DyingStatus: int32(pd.DyingStatus),
	}
}

type Player struct {
	mu             sync.Mutex
	id             string
	dx             float32
	dy             float32
	angle          float32 // degree
	speed          float32
	isMoveing      bool
	isDashing      bool
	dashCoolTime   int
	health         int
	MagicType      int
	DyingStatus    *PlayerDead
	collisionList  []int
	physicalObject physical.Physical
}

func NewPlayer(id string, x float32, y float32) *Player {

	return &Player{
		sync.Mutex{},
		id,
		0,
		0,
		0,
		PLAYER_SPEED,
		false,
		false,
		0,
		100,
		BULLET_NONE,
		NewPlayerDead("", id, DYING_NONE),
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
		player.physicalObject.Move(player.dx, player.dy)
		if player.isDashing {
			player.dashCoolTime--
			if player.dashCoolTime == 0 {
				player.mu.Lock()
				player.isDashing = false
				player.speed = PLAYER_SPEED
				player.mu.Unlock()
			}
		}
	}
	player.mu.Unlock()
}

func (player *Player) IsValid() bool {
	return player.health > 0
}

func (player *Player) DirChange(angle float64, isMoved bool) {
	player.mu.Lock()
	player.dx = player.speed * float32(math.Sin(angle*(math.Pi/180)))
	player.dy = player.speed * float32(math.Cos(angle*(math.Pi/180))) * -1
	player.isMoveing = isMoved
	player.mu.Unlock()
}

func (player *Player) AngleChange(angle float32) {
	player.angle = angle
}

func (player *Player) MakeSendingData() *message.PlayerState {
	return &message.PlayerState{
		Id:           player.id,
		X:            player.physicalObject.GetX(),
		Y:            player.physicalObject.GetY(),
		Health:       int32(player.health),
		MagicType:    int32(player.MagicType),
		Angle:        float32(player.angle),
		DashCoolTime: int32(player.dashCoolTime),
	}
}

func (player *Player) HeatedBullet(bullet *Bullet) {
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
		player.dashCoolTime = DASH_COOLTIME
		player.mu.Unlock()
	}
}

func (player *Player) GetPhysical() *physical.Physical {
	return &player.physicalObject
}

func (player *Player) addCollider(objectType int, effect func(obj Object)) {

}
