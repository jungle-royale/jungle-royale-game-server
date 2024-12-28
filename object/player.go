package object

import (
	"jungle-royale/cons"
	"jungle-royale/message"
	"jungle-royale/object/physical"
	"log"
	"math"
	"sync"
	"time"
)

const PLAYER_SPEED = 0.3
const DASH_SPEED = 1.2
const DASH_TICK = 5
const DASH_COOLTIME = 120 // 2sec
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

func (player *Player) CalcCollision() {

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
		Id: player.id,
		X:  player.physicalObject.GetX(),
		Y:  player.physicalObject.GetY(),
	}
}

func (player *Player) HeatedBullet() {
	player.mu.Lock()
	player.health -= BULLET_DAMAGE
	player.mu.Unlock()
}

func (player *Player) DoDash() {
	if !player.isDashing && player.dashCoolTime < 0 {
		log.Printf("dash")
		player.mu.Lock()
		player.isDashing = true
		player.speed = DASH_SPEED
		player.mu.Unlock()
		time.AfterFunc(cons.CalcLoopInterval*DASH_TICK*time.Millisecond, func() {
			log.Printf("dash end")
			player.mu.Lock()
			player.isDashing = false
			player.speed = PLAYER_SPEED
			player.dashCoolTime = DASH_COOLTIME
			player.mu.Unlock()
		})
	}
}
