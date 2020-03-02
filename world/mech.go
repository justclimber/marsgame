package world

import "aakimov/marsgame/physics"

type Mech struct {
	Object
	cannon         *Cannon
	throttle       float64
	rotateThrottle float64
	generator      Generator
}

type Cannon struct {
	shoot          Shoot
	rotateThrottle float64
	angle          float64
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

func NewMech(x, y float64) *Mech {
	return &Mech{
		Object: NewObject(
			1,
			TypePlayer,
			physics.Point{x, y},
			100,
			0,
			0,
			0,
			1000,
		),
		cannon: &Cannon{
			rotateThrottle: 0,
			angle:          0,
		},
		throttle:       0,
		rotateThrottle: 0,
		generator: Generator{
			efficiency:            900,
			rateMs:                100,
			value:                 20000,
			maxValue:              40000,
			xelons:                1000,
			terminateCh:           make(chan bool, 1),
			terminateThrottlingCh: make(chan bool, 1),
		},
	}
}
