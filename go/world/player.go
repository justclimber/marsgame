package world

import (
	"aakimov/marsgame/go/server"
	"log"
)

type Player struct {
	id          string
	client      *server.Client
	mech        Mech
	mainProgram *Code
}

func NewPlayer(id string, client *server.Client, mech Mech) *Player {
	return &Player{
		id:     id,
		client: client,
		mech:   mech,
		mainProgram: &Code{
			id: "main",
		},
	}
}

func (p *Player) setBaseParams() {
	p.mech.Throttle = 1
}

func (p *Player) saveAstCode(sourceCode string) {
	astCode, _ := ParseSourceCode(sourceCode)
	p.mainProgram.saveAst(astCode)
	log.Println("New ast code saved")
}
