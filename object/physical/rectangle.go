package physical

import (
	"log"
	"math"
)

type Rectangle struct {
	X      float64
	Y      float64 // left top
	Dx     float64
	Dy     float64
	Width  float64
	Length float64
}

func NewRectangle(x, y, width, length float64) *Rectangle {
	return &Rectangle{
		x,
		y,
		0,
		0,
		width,
		length,
	}
}

func (rectangle *Rectangle) pointInRectangle(x, y float64) bool {
	if rectangle.X < x && x < rectangle.X+rectangle.Width &&
		rectangle.Y < y && y < rectangle.Y+rectangle.Length {
		return true
	} else {
		return false
	}
}

func (obj *Rectangle) IsCollide(opponent *Physical) bool { // opponent is pointer

	switch opp := (*opponent).(type) {

	case *Circle:

		closestX := math.Max(float64(obj.X), math.Min(float64(opp.X), float64(obj.X+obj.Width)))
		closestY := math.Max(float64(obj.Y), math.Min(float64(opp.Y), float64(obj.Y+obj.Length)))
		distanceX := float64(opp.X) - closestX
		distanceY := float64(opp.Y) - closestY
		distanceSquared := distanceX*distanceX + distanceY*distanceY

		if distanceSquared <= float64(opp.Radius*opp.Radius) {
			return true
		} else {
			return false
		}

	case *Rectangle:
		if obj.pointInRectangle(opp.X, opp.Y) ||
			obj.pointInRectangle(opp.X+opp.Width, opp.Y+opp.Length) {
			return true
		} else {
			return false
		}

	default:
		log.Printf("err: type is unmatched")
		return false
	}
}

func (rect *Rectangle) Move() {
	rect.X += rect.Dx
	rect.Y += rect.Dy
}

func (rect *Rectangle) GetX() float64 {
	return rect.X
}

func (rect *Rectangle) GetY() float64 {
	return rect.Y
}

func (rect *Rectangle) SetCoord(x, y float64) {
	rect.X = x
	rect.Y = y
}

func (rect *Rectangle) IsInRectangle(x, y float64) bool {
	if rect.X <= x && x <= rect.X+rect.Width &&
		rect.Y <= y && y <= rect.Y+rect.Length {
		return true
	} else {
		return false
	}
}

func (r *Rectangle) GetDx() float64 {
	return r.Dx
}

func (r *Rectangle) GetDy() float64 {
	return r.Dy
}

func (r *Rectangle) SetDir(dx, dy float64) {
	r.Dx = dx
	r.Dy = dy
}

func (r *Rectangle) CollideRelocate(obj *Physical) {

}
