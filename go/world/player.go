package world

import (
	"aakimov/marsgame/go/code"
	"aakimov/marsgame/go/server"
)

type Player struct {
	id          string
	client      *server.Client
	mech        Mech
	mainProgram *code.Code
}

func NewPlayer(id string, client *server.Client, mech Mech) *Player {
	return &Player{
		id:     id,
		client: client,
		mech:   mech,
	}
}

func (p *Player) setBaseParams() {
	p.mech.Throttle = 1
}
