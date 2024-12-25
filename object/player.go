package object

import (
	"jungle-royale/message"
	"math"
	"sync"
)

const PLAYER_SPEED = 5
const EPSILON = 1e-9

type Player struct {
	lock    sync.Mutex
	id      string
	x       float32
	y       float32
	radious float32
	angle   float64 // degree
	isMoved bool
}

func NewPlayer(id string, x float32, y float32) *Player {
	return &Player{sync.Mutex{}, id, x, y, 30, 0, false}
}

func (player *Player) Move() {
	if player.isMoved {
		player.lock.Lock()
		player.x += PLAYER_SPEED * float32(math.Sin(player.angle*(math.Pi/180)))
		player.y -= PLAYER_SPEED * float32(math.Cos(player.angle*(math.Pi/180)))
		player.lock.Unlock()
	}
}

func (player *Player) DirChange(angle float64, isMoved bool) {
	player.lock.Lock()
	player.angle = angle
	player.isMoved = isMoved
	player.lock.Unlock()
}

func (player *Player) MakeSendingData() *message.Player {
	return &message.Player{
		Id: player.id,
		X:  player.x,
		Y:  player.y,
	}
}
