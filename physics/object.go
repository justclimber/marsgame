package physics

import (
	"sync"
)

type Obj struct {
	Id   uint32
	Type string
	sync.Mutex
	Pos             Point
	CollisionRadius int
	Angle           float64
	Speed           float64
	AngleSpeed      float64
	Weight          float64
	Velocity        *Vector
	Direction       *Vector
}
