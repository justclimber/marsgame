package world

import (
	"aakimov/marsgame/helpers"
	"aakimov/marsgame/server"
	"log"
	"time"
)

type Player struct {
	id          string
	client      *server.Client
	mech        *Mech
	mainProgram *Code
}

func NewPlayer(id string, client *server.Client, w *World, runSpeedMs time.Duration) *Player {
	mech := NewMech(1000, 1000)
	player := &Player{
		id:          id,
		client:      client,
		mech:        mech,
		mainProgram: NewCode("main", w, mech, runSpeedMs),
	}
	return player
}

func (p *Player) setBaseParams() {
	p.mech.Throttle = 1
}

func (p *Player) saveAstCode(sourceCode string) {
	p.mainProgram.saveCode(sourceCode)
}

func (p *Player) listen() {
	log.Printf("Player [%s] listening started", p.id)
	for {
		select {
		case codeOutputs := <-p.mainProgram.outputCh:
			log.Printf("Write code run result for player [%s]: mThr: %f, mrThr: %f, crThr: %f, shoot: %f",
				p.id, codeOutputs.MThrottle, codeOutputs.RThrottle, codeOutputs.CRThrottle, codeOutputs.Shoot)
			p.mech.mu.Lock()
			p.mech.Throttle = codeOutputs.MThrottle
			p.mech.RotateThrottle = codeOutputs.RThrottle
			p.mech.Cannon.RotateThrottle = codeOutputs.CRThrottle
			if codeOutputs.Shoot != 0 {
				p.mech.Cannon.shoot.state = WillShoot
				p.mech.Cannon.shoot.delay = int(codeOutputs.Shoot * 1000)
			}
			p.mech.mu.Unlock()
		case codeError := <-p.mainProgram.errorCh:
			p.client.PackAndSendCommand("codeError", codeError)
		}
	}
}

const nearTimeDelta = 50

func areTimeIdNearlySameOrGrater(t1, t2 int64) bool {
	return t1 > t2 || helpers.AbsInt64(t1-t2) < nearTimeDelta
}

func (p *Player) run(world *World) *ChangeByObject {
	mech := p.mech
	changeByObject := ChangeByObject{
		ObjType: TypePlayer,
		ObjId:   p.id,
	}
	mech.mu.Lock()
	if mech.Cannon.shoot.state == WillShoot {
		mech.Cannon.shoot.state = Planned
		mech.Cannon.shoot.willShootAt = world.timeId + int64(mech.Cannon.shoot.delay)
	}
	if mech.Cannon.shoot.state == Planned && areTimeIdNearlySameOrGrater(world.timeId, mech.Cannon.shoot.willShootAt) {
		mech.Cannon.shoot.state = None
		p.shoot(world)
	}
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

func (p *Player) shoot(world *World) {
	log.Println("Shoooooot!!!!")
	missileAngle := p.mech.Cannon.Angle + p.mech.Angle
	missilePos := p.mech.Pos

	//move missile a bit of forward far away from mech center
	missilePos.MoveForward(missileAngle, 100.)
	world.newObjectsCh <- &Missile{
		Object: Object{
			Type:            TypeMissile,
			Speed:           MissileSpeed,
			Pos:             missilePos,
			Angle:           missileAngle,
			CollisionRadius: 20,
		},
	}
}
