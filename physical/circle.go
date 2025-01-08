package physical

import (
	"jungle-royale/util"
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
		distanceSquared := math.Pow(obj.X-opp.X, 2) + math.Pow(obj.Y-opp.Y, 2)
		radiousSquared := math.Pow(obj.Radius+opp.Radius, 2)
		if radiousSquared > distanceSquared {
			return true
		} else {
			return false
		}

	case *Rectangle:

		closestX := math.Max(opp.X, math.Min(obj.X, opp.X+opp.Width))
		closestY := math.Max(opp.Y, math.Min(obj.Y, opp.Y+opp.Length))
		distanceX := obj.X - closestX
		distanceY := obj.Y - closestY
		distanceSquared := distanceX*distanceX + distanceY*distanceY

		if distanceSquared <= obj.Radius*obj.Radius {
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

func clamp(value, minVal, maxVal float64) float64 {
	if value < minVal {
		return minVal
	}
	if value > maxVal {
		return maxVal
	}
	return value
}

func (c *Circle) CollideRelocate(obj *Physical) {

	switch opp := (*obj).(type) {

	case *Circle:

		dx := opp.X - c.X
		dy := opp.Y - c.Y
		d := math.Hypot(dx, dy)
		if d == 0 {
			return
		}
		nx := dx / d
		ny := dy / d

		pen := (c.Radius + opp.Radius) - d
		if pen > 0 {
			c.X -= nx * pen
			c.Y -= ny * pen
		}

	case *Rectangle:

		closestX := clamp(c.X, opp.X, opp.X+opp.Width)
		closestY := clamp(c.Y, opp.Y, opp.Y+opp.Length)

		dx := c.X - closestX
		dy := c.Y - closestY
		distSq := dx*dx + dy*dy

		radiusSq := c.Radius * c.Radius
		if distSq <= radiusSq {

			dist := math.Sqrt(distSq)
			if dist == 0 {
				c.Y += c.Radius
				return
			}

			pen := c.Radius - dist

			nx := dx / dist
			ny := dy / dist

			c.X += nx * pen
			c.Y += ny * pen
		}

	default:
		log.Printf("err: type is unmatched")
		return
	}
}

func (c *Circle) GetBoundCoordSet() *util.Set[Coord] {
	ret := util.NewSet[Coord]()
	ret.Add(Coord{c.X - c.Radius, c.Y})
	ret.Add(Coord{c.X + c.Radius, c.Y})
	ret.Add(Coord{c.X, c.Y + c.Radius})
	ret.Add(Coord{c.X, c.Y - c.Radius})
	return ret
}
