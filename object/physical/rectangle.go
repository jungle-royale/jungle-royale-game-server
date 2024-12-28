package physical

import (
	"log"
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
		if opp.pointInCircle(obj.X, obj.Y) ||
			opp.pointInCircle(obj.X+obj.Width, obj.Y) ||
			opp.pointInCircle(obj.X, obj.Y+obj.Length) ||
			opp.pointInCircle(obj.X+obj.Width, obj.Y+obj.Length) {
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
