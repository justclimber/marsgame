package physics

import (
	"sync"
)

type Obj struct {
	Id   int
	Type string
	sync.Mutex
	Pos             Point
	CollisionRadius int
	Angle           float64
	Speed           float64
	AngleSpeed      float64
	Weight          float64
	MoveDone        *Vector
	Velocity        *Vector
	Direction       *Vector
}
