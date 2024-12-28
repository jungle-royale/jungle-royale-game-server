package physical

import (
	"log"
	"math"
)

type Circle struct {
	X       float32
	Y       float32
	Radious float32
}

func NewCircle(x float32, y float32, radious float32) *Circle {
	return &Circle{
		x,
		y,
		radious,
	}
}

func (circle *Circle) pointInCircle(x float32, y float32) bool {
	distanceSquared := math.Pow(float64(circle.X-x), 2) + math.Pow(float64(circle.Y-y), 2)
	radiousSquared := math.Pow(float64(circle.Radious), 2)
	if radiousSquared > distanceSquared {
		return true
	} else {
		return false
	}
}

func (obj *Circle) IsCollide(opponent *Physical) bool { // opponent is pointer

	switch opp := (*opponent).(type) {

	case *Circle:
		distanceSquared := math.Pow(float64(obj.X-opp.X), 2) + math.Pow(float64(obj.Y-opp.Y), 2)
		radiousSquared := math.Pow(float64(obj.Radious+opp.Radious), 2)
		if radiousSquared > distanceSquared {
			return true
		} else {
			return false
		}

	case *Rectangle:
		if obj.pointInCircle(opp.X, opp.Y) ||
			obj.pointInCircle(opp.X+opp.Width, opp.Y) ||
			obj.pointInCircle(opp.X, opp.Y+opp.Length) ||
			obj.pointInCircle(opp.X+opp.Width, opp.Y+opp.Length) {
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
