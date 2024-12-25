package object

type Player struct {
	id      string
	x       float32
	y       float32
	radious float32
	dx      float32
	dy      float32
}

func NewPlayer(id string, x float32, y float32) *Player {
	return &Player{id, x, y, 30, 0, 0}
}

func (player *Player) Move() {
	player.x += player.dx
	player.y += player.dy
}

func (player *Player) DirChange(dx float32, dy float32) {
	player.dx = dx
	player.dy = dy
}
