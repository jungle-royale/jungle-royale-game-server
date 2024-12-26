package object

import (
	"jungle-royale/message"
	"math"
)

const PLAYER_SPEED = 5
const EPSILON = 1e-9

type Player struct {
	id      string
	x       float32
	y       float32
	radious float32
	angle   float64 // degree
	isMoved bool
	health  int
}

func NewPlayer(id string, x float32, y float32) *Player {
	return &Player{id, x, y, 30, 0, false, 100}
}

func (player *Player) Move() {
	if player.isMoved {
		player.x += PLAYER_SPEED * float32(math.Sin(player.angle*(math.Pi/180)))
		player.y -= PLAYER_SPEED * float32(math.Cos(player.angle*(math.Pi/180)))
	}
}

func (player *Player) IsValid() bool {
	if player.health <= 0 {
		return false
	}
	return true
}

func (player *Player) DirChange(angle float64, isMoved bool) {
	player.angle = angle
	player.isMoved = isMoved
}

func (player *Player) MakeSendingData() *message.Player {
	return &message.Player{
		Id: player.id,
		X:  player.x,
		Y:  player.y,
	}
}

func (player *Player) HeatedBullet() {
	player.health -= BULLET_DAMAGE
}
