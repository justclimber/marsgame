package world

type Mech struct {
	Object
	Throttle       float64
	RotateThrottle float64
}

func NewMech() Mech {
	return Mech{
		Object:         Object{},
		Throttle:       0,
		RotateThrottle: 0,
	}
}
