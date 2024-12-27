package object

import (
	"jungle-royale/cons"
	"jungle-royale/message"
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
	mu           sync.Mutex
	id           string
	x            float32
	y            float32
	radious      float32
	angle        float64 // degree
	speed        float32
	isMoveing    bool
	isDashing    bool
	dashCoolTime int
	health       int
}

func NewPlayer(id string, x float32, y float32) *Player {
	return &Player{
		sync.Mutex{},
		id,
		x,
		y,
		PLAYER_RADIOUS,
		0,
		PLAYER_SPEED,
		false,
		false,
		0,
		100,
	}
}

func (player *Player) SetLocation(x float32, y float32) {
	player.mu.Lock()
	player.x = x
	player.y = y
	player.mu.Unlock()
}

func (player *Player) Move() {
	player.mu.Lock()
	if player.isMoveing {
		player.x += player.speed * float32(math.Sin(player.angle*(math.Pi/180)))
		player.y -= player.speed * float32(math.Cos(player.angle*(math.Pi/180)))
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
		Id: player.id,
		X:  player.x,
		Y:  player.y,
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
