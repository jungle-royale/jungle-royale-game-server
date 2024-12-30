package physical

import (
	"log"
	"math"
)

type Rectangle struct {
	X      float32
	Y      float32 // left top
	Width  float32
	Length float32
}

func NewRectangle(x float32, y float32, width float32, length float32) *Rectangle {
	return &Rectangle{
		x,
		y,
		width,
		length,
	}
}

func (rectangle *Rectangle) pointInRectangle(x float32, y float32) bool {
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

func (rect *Rectangle) Move(dx float32, dy float32) {
	rect.X += dx
	rect.Y += dy
}

func (rect *Rectangle) GetX() float32 {
	return rect.X
}

func (rect *Rectangle) GetY() float32 {
	return rect.Y
}

func (rect *Rectangle) SetCoord(x float32, y float32) {
	rect.X = x
	rect.Y = y
}

func (rect *Rectangle) IsInRectangle(x, y float32) bool {
	if rect.X < x && x < rect.X+rect.Width &&
		rect.Y < y && y < rect.Y+rect.Length {
		return true
	} else {
		return false
	}
}
