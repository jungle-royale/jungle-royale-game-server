package object

import (
	"jungle-royale/message"
	"sync"
)

type Player struct {
	lock    sync.Mutex
	id      string
	x       float32
	y       float32
	radious float32
	dx      float32
	dy      float32
}

func NewPlayer(id string, x float32, y float32) *Player {
	return &Player{sync.Mutex{}, id, x, y, 30, 0, 0}
}

func (player *Player) Move() {
	player.lock.Lock()
	player.x += player.dx
	player.y += player.dy
	player.lock.Unlock()
}

func (player *Player) DirChange(dx float32, dy float32) {
	player.lock.Lock()
	player.dx = dx
	player.dy = dy
	player.lock.Unlock()
}

func (player *Player) MakeSendingPlayerData() *message.Player {
	return &message.Player{
		Id: player.id,
		X:  player.x,
		Y:  player.y,
	}
}
