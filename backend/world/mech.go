package world

import "sync"

type Mech struct {
	mu sync.Mutex
	Object
	Cannon         *Cannon
	Throttle       float64
	RotateThrottle float64
}

type Cannon struct {
	RotateThrottle float64
	Angle          float64
}

func NewMech() *Mech {
	return &Mech{
		Object: Object{},
		Cannon: &Cannon{
			RotateThrottle: 0,
			Angle:          0,
		},
		Throttle:       0,
		RotateThrottle: 0,
	}
}
