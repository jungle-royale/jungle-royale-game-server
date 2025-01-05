package physical

// nonmoving object
type Physical interface {
	IsCollide(opponent *Physical) bool
	Move()
	GetX() float64
	GetY() float64
	GetDx() float64
	GetDy() float64
	SetCoord(x, y float64)
	CollideRelocate(obj *Physical)
	SetDir(dx, dy float64)
}
