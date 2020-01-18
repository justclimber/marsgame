package world

import (
	"aakimov/marsgame/backend/server"
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
		id:          id,
		client:      client,
		mech:        mech,
		mainProgram: NewCode("main", w, mech),
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
			log.Printf("Write code run result for player [%s]: mThr: %f, mrThr: %f, crThr: %f",
				p.id, codeOutputs.MThrottle, codeOutputs.RThrottle, codeOutputs.СRThrottle)
			p.mech.mu.Lock()
			p.mech.Throttle = codeOutputs.MThrottle
			p.mech.RotateThrottle = codeOutputs.RThrottle
			p.mech.Cannon.RotateThrottle = codeOutputs.СRThrottle
			p.mech.mu.Unlock()
		case codeError := <-p.mainProgram.errorCh:
			p.client.PackAndSendCommand("error", codeError)
		default:
			// noop
		}
	}
}

func (p *Player) run(world *World) *ChangeByObject {
	mech := p.mech
	changeByObject := ChangeByObject{
		ObjType: TypePlayer,
		ObjId:   p.id,
	}
	mech.mu.Lock()
	if mech.RotateThrottle != 0 {
		mech.Object.Angle += mech.RotateThrottle * MaxRotationValue
		newAngle := mech.Object.Angle
		changeByObject.Angle = &newAngle
	}
	if mech.Throttle != 0 {
		length := mech.Throttle * MaxMovingLength
		mech.Object.Pos.MoveForward(mech.Object.Angle, length)
		newPos := mech.Object.Pos
		changeByObject.Pos = &newPos
		changeByObject.length = &length
	}
	if mech.Cannon.RotateThrottle != 0 {
		mech.Cannon.Angle += mech.Cannon.RotateThrottle * MaxCannonRotationValue
		newCannonAngle := mech.Cannon.Angle
		changeByObject.CannonAngle = &newCannonAngle
	}
	mech.mu.Unlock()

	if mech.RotateThrottle != 0 || mech.Throttle != 0 || mech.Cannon.RotateThrottle != 0 {
		return &changeByObject
	}
	return nil
}
