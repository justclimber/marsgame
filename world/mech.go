package world

import "aakimov/marsgame/physics"

type Mech struct {
	Object
	cannon    *Cannon
	commands  *Commands
	generator Generator
}

type Commands struct {
	cannon *CannonCommands
	mech   *MechCommands
}

func NewEmptyCommands() *Commands {
	return &Commands{
		cannon: &CannonCommands{&Shoot{}, 0},
		mech:   &MechCommands{},
	}
}

type MechCommands struct {
	move   float64
	rotate float64
}

type CannonCommands struct {
	shoot  *Shoot
	rotate float64
}

type Cannon struct {
	shoot *Shoot
	angle float64
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
			physics.Point{X: x, Y: y},
			100,
			0,
			0,
			0,
			1000,
		),
		cannon:   &Cannon{angle: 0, shoot: &Shoot{}},
		commands: NewEmptyCommands(),
		generator: Generator{
			efficiency:            900,
			rateMs:                100 / TimeMultiplicator,
			value:                 20000,
			maxValue:              40000,
			xelons:                1000,
			terminateCh:           make(chan bool, 1),
			terminateThrottlingCh: make(chan bool, 1),
		},
	}
}
