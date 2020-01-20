package world

import (
	"aakimov/marsgame/physics"
	"sync"
)

type Mech struct {
	mu sync.Mutex
	Object
	Cannon         *Cannon
	Throttle       float64
	RotateThrottle float64
}

type Cannon struct {
	shoot          Shoot
	RotateThrottle float64
	Angle          float64
}

const (
	None = iota
	WillShoot
	Planned
)

type Shoot struct {
	state       int
	delay       int
	willShootAt int64
}

type Missile struct {
	Object
}

func NewMech(x, y float64) *Mech {
	return &Mech{
		Object: Object{Pos: physics.Point{
			X: x,
			Y: y,
		}},
		Cannon: &Cannon{
			RotateThrottle: 0,
			Angle:          0,
		},
		Throttle:       0,
		RotateThrottle: 0,
	}
}
