package world

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

type Missile struct {
	Object
}

func NewMech(x, y float64) *Mech {
	return &Mech{
		Object: NewObject(
			1,
			"player",
			Point{x, y},
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
			increment: 200,
			rateMs:    200,
			value:     2000,
			maxValue:  20000,
		},
	}
}
