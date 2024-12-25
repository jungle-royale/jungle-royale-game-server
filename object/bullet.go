package object

import "sync"

const RANGE = 100

type Bullet struct {
	lock     sync.Mutex
	playerId string
	dx       float32
	dy       float32
	endX     float32
	endY     float32
}
