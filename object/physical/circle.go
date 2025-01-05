package physical

import (
	"log"
	"math"
)

type Circle struct {
	X      float64
	Y      float64
	Dx     float64
	Dy     float64
	Radius float64
}

func NewCircle(x, y, Radius float64) *Circle {
	return &Circle{
		x,
		y,
		0,
		0,
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

func (circle *Circle) Move() {
	circle.X += circle.Dx
	circle.Y += circle.Dy
}

func (circle *Circle) GetX() float64 {
	return circle.X
}

func (circle *Circle) GetY() float64 {
	return circle.Y
}

func (circle *Circle) SetCoord(x, y float64) {
	circle.X = x
	circle.Y = y
}

func (c *Circle) GetDx() float64 {
	return c.Dx
}

func (c *Circle) GetDy() float64 {
	return c.Dy
}

func (c *Circle) SetDir(dx, dy float64) {
	c.Dx = dx
	c.Dy = dy
}

func (c *Circle) CollideRelocate(obj *Physical) {

	var t, ox, oy, th, fx, fy float64

	switch opp := (*obj).(type) {

	case *Circle:
		xx := c.X - opp.X
		yy := c.Y - opp.Y
		rr := c.Radius + opp.Radius
		a := c.Dx*c.Dx + c.Dy*c.Dy
		b := xx*c.Dx + yy*c.Dy
		cC := xx*xx + yy*yy - rr*rr
		if a == 0 {
			return
		}
		disc := b*b - a*cC
		if disc < 0 {
			return
		}
		t = (2*b + math.Sqrt(math.Abs(4*b*b-4*a*cC))) / (2 * a)

		ox = c.X - c.Dx*t
		oy = c.Y - c.Dy*t

		deltaY := oy - opp.Y
		deltaX := ox - opp.X
		th = math.Atan2(deltaY, deltaX)

	case *Rectangle:

	default:
		log.Printf("err: type is unmatched")
		return
	}

	va := math.Sin(-th)*c.Dy*t + math.Cos(th)*c.Dx*t

	ax := va * math.Sin(th)
	ay := va * math.Cos(th)

	fx = ox + ax
	fy = oy + ay

	c.SetCoord(fx, fy)
}
