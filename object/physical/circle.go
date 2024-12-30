package physical

import (
	"log"
	"math"
)

type Circle struct {
	X      float32
	Y      float32
	Radius float32
}

func NewCircle(x float32, y float32, Radius float32) *Circle {
	return &Circle{
		x,
		y,
		Radius,
	}
}

func (obj *Circle) IsCollide(opponent *Physical) bool { // opponent is pointer

	switch opp := (*opponent).(type) {

	case *Circle:
		distanceSquared := math.Pow(float64(obj.X-opp.X), 2) + math.Pow(float64(obj.Y-opp.Y), 2)
		radiousSquared := math.Pow(float64(obj.Radius+opp.Radius), 2)
		if radiousSquared > distanceSquared {
			return true
		} else {
			return false
		}

	case *Rectangle:

		closestX := math.Max(float64(opp.X), math.Min(float64(obj.X), float64(opp.X+opp.Width)))
		closestY := math.Max(float64(opp.Y), math.Min(float64(obj.Y), float64(opp.Y+opp.Length)))
		distanceX := float64(obj.X) - closestX
		distanceY := float64(obj.Y) - closestY
		distanceSquared := distanceX*distanceX + distanceY*distanceY

		if distanceSquared <= float64(obj.Radius*obj.Radius) {
			return true
		} else {
			return false
		}

	default:
		log.Printf("err: type is unmatched")
		return false
	}
}

func (circle *Circle) Move(dx float32, dy float32) {
	circle.X += dx
	circle.Y += dy
}

func (circle *Circle) GetX() float32 {
	return circle.X
}

func (circle *Circle) GetY() float32 {
	return circle.Y
}

func (circle *Circle) SetCoord(x float32, y float32) {
	circle.X = x
	circle.Y = y
}
