package world

import (
	"aakimov/marsgame/go/server"
	"log"
)

type Player struct {
	id          string
	client      *server.Client
	mech        *Mech
	mainProgram *Code
}

func NewPlayer(id string, client *server.Client, mech *Mech, w *World) *Player {
	player := &Player{
		id:     id,
		client: client,
		mech:   mech,
		mainProgram: &Code{
			id:       "main",
			outputCh: make(chan *MechOutputVars),
			worldP:   w,
			mechP:    mech,
		},
	}
	return player
}

func (p *Player) setBaseParams() {
	p.mech.Throttle = 1
}

func (p *Player) saveAstCode(sourceCode string) {
	astCode, _ := ParseSourceCode(sourceCode)
	p.mainProgram.saveAst(astCode)
	log.Println("New ast code saved")
}

func (p *Player) listen() {
	log.Printf("Player [%s] listening started", p.id)
	for {
		select {
		case codeOutputs := <-p.mainProgram.outputCh:
			log.Printf("Write code run result for player [%s]: mThrottle: %f, rThrottle: %f",
				p.id, codeOutputs.MThrottle, codeOutputs.RThrottle)
			p.mech.mu.Lock()
			p.mech.Throttle = codeOutputs.MThrottle
			p.mech.RotateThrottle = codeOutputs.RThrottle
			p.mech.mu.Unlock()
		default:
			// noop
		}
	}
}
