package world

import (
	"aakimov/marsgame/helpers"
	"aakimov/marsgame/server"
	"aakimov/marslang/ast"
	"aakimov/marslang/lexer"
	"aakimov/marslang/parser"
	"log"
	"math"
	"time"
)

type Player struct {
	id          string
	world       *World
	client      *server.Client
	mech        *Mech
	mainProgram *Code
	runSpeedMs  time.Duration
	outputCh    chan *MechOutputVars
	io4ClientCh chan *IO4Client
	codeSaveCh  chan *ast.StatementsBlock
	flowCh      chan ProgramState
	errorCh     chan *Error
}

func NewPlayer(id string, client *server.Client, w *World, runSpeedMs time.Duration) *Player {
	mech := NewMech(1000, 1000)
	player := &Player{
		id:          id,
		world:       w,
		client:      client,
		mech:        mech,
		mainProgram: NewCode("main", w, mech, runSpeedMs),
		runSpeedMs:  runSpeedMs,
		outputCh:    make(chan *MechOutputVars),
		codeSaveCh:  make(chan *ast.StatementsBlock),
		io4ClientCh: make(chan *IO4Client),
		flowCh:      make(chan ProgramState),
		errorCh:     make(chan *Error),
	}
	return player
}

func (p *Player) setBaseParams() {
	p.mech.throttle = 1
}

func (p *Player) saveAstCode(sourceCode string) {
	l := lexer.New(sourceCode)
	pr, err := parser.New(l)
	if err != nil {
		p.errorCh <- &Error{
			ErrorType: Lexing,
			Message:   err.Error(),
		}
		log.Printf("Lexing error: %s", err.Error())
		return
	}
	astProgram, err := pr.Parse()
	if err != nil {
		p.errorCh <- &Error{
			ErrorType: Parsing,
			Message:   err.Error(),
		}
		log.Printf("Parsing error: %s", err.Error())
		return
	}
	log.Println("Code parsed")
	p.codeSaveCh <- astProgram
}

func (p *Player) listen() {
	log.Printf("Player [%s] listening started", p.id)
	for {
		select {
		case codeOutputs := <-p.outputCh:
			log.Printf("Write code runProgram result for player [%s]: mThr: %f, mrThr: %f, crThr: %f, shoot: %f",
				p.id, codeOutputs.MThrottle, codeOutputs.RThrottle, codeOutputs.CRThrottle, codeOutputs.Shoot)
			p.mech.Lock()
			p.mech.throttle = codeOutputs.MThrottle
			p.mech.rotateThrottle = codeOutputs.RThrottle
			p.mech.cannon.rotateThrottle = codeOutputs.CRThrottle
			if codeOutputs.Shoot != 0 {
				p.mech.cannon.shoot.state = WillShoot
				p.mech.cannon.shoot.delay = int(codeOutputs.Shoot * 1000)
			}
			p.mech.Unlock()
		case codeError := <-p.errorCh:
			p.client.PackAndSendCommand("codeError", codeError)
		case io4Client := <-p.io4ClientCh:
			p.client.PackAndSendCommand("codeInputOutput", io4Client)
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
	mech.Lock()
	defer mech.Unlock()

	if mech.cannon.shoot.state == WillShoot {
		mech.cannon.shoot.state = Planned
		mech.cannon.shoot.willShootAt = world.timeId + int64(mech.cannon.shoot.delay)
	}
	if mech.cannon.shoot.state == Planned && areTimeIdNearlySameOrGrater(world.timeId, mech.cannon.shoot.willShootAt) {
		mech.cannon.shoot.state = None
		p.shoot(world)
	}
	if mech.rotateThrottle != 0 {
		mech.Object.Angle += mech.rotateThrottle * MaxRotationValue
		if mech.Object.Angle > 2*math.Pi {
			mech.Object.Angle = mech.Object.Angle - 2*math.Pi
		} else if mech.Object.Angle < 0 {
			mech.Object.Angle = 2*math.Pi + mech.Object.Angle
		}
		newAngle := mech.Object.Angle
		changeByObject.Angle = &newAngle
	}
	if mech.throttle != 0 {
		length := mech.throttle * MaxMovingLength
		mech.Object.Pos.MoveForward(mech.Object.Angle, length)
		newPos := mech.Object.Pos
		changeByObject.Pos = &newPos
		changeByObject.length = &length
	}
	if mech.cannon.rotateThrottle != 0 {
		mech.cannon.angle += mech.cannon.rotateThrottle * MaxCannonRotationValue
		newCannonAngle := mech.cannon.angle
		changeByObject.CannonAngle = &newCannonAngle
	}

	if mech.rotateThrottle != 0 || mech.throttle != 0 || mech.cannon.rotateThrottle != 0 {
		return &changeByObject
	}
	return nil
}

func (p *Player) shoot(world *World) {
	missileAngle := p.mech.cannon.angle + p.mech.Angle
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
