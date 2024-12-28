package physical

// nonmoving object
type Physical interface {
	IsCollide(opponent *Physical) bool
	Move(dx float32, dy float32)
	GetX() float32
	GetY() float32
	SetCoord(x float32, y float32)
}
