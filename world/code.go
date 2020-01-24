package world

import (
	"aakimov/marsgame/server"
	"aakimov/marslang/ast"
	"aakimov/marslang/interpereter"
	"aakimov/marslang/lexer"
	"aakimov/marslang/object"
	"aakimov/marslang/parser"

	"log"
	"sync"
	"time"
)

type ProgramState int

const (
	Stopped ProgramState = iota
	Running
)

type Code struct {
	id         string
	state      ProgramState
	mu         sync.Mutex
	astProgram *ast.StatementsBlock
	worldP     *World
	mechP      *Mech
	outputCh   chan *MechOutputVars
	codeSaveCh chan *ast.StatementsBlock
	flowCh     chan ProgramState
	errorCh    chan *Error
	runSpeedMs time.Duration
}

func NewCode(id string, world *World, mech *Mech, runSpeedMs time.Duration) *Code {
	return &Code{
		id:         "main",
		worldP:     world,
		mechP:      mech,
		outputCh:   make(chan *MechOutputVars),
		codeSaveCh: make(chan *ast.StatementsBlock),
		flowCh:     make(chan ProgramState),
		errorCh:    make(chan *Error),
		runSpeedMs: runSpeedMs,
	}
}

type MechOutputVars struct {
	MThrottle  float64
	RThrottle  float64
	CRThrottle float64
	Shoot      float64
}

type ErrorType int

const (
	Lexing ErrorType = iota
	Parsing
	Runtime
)

type Error struct {
	ErrorType ErrorType `json:"errorType"`
	Message   string    `json:"message"`
}

func newMechOutputVarsFromEnv(env *object.Environment) *MechOutputVars {
	return &MechOutputVars{
		MThrottle:  getFloatVarFromEnv("mThr", env),
		RThrottle:  getFloatVarFromEnv("mrThr", env),
		CRThrottle: getFloatVarFromEnv("crThr", env),
		Shoot:      getFloatVarFromEnv("shoot", env),
	}
}

func getFloatVarFromEnv(varName string, env *object.Environment) float64 {
	envObj, ok := env.Get(varName)
	if !ok {
		return 0
	}
	objFloat, ok := envObj.(*object.Float)
	if !ok {
		log.Fatalf("%s should be type float, %s given", varName, envObj.Type())
	}
	return objFloat.Value
}

func (c *Code) loadMechVarsIntoEnv(env *object.Environment) {
	s := make(map[string]interface{})
	s["x"] = c.mechP.Pos.x
	s["y"] = c.mechP.Pos.y
	s["angle"] = c.mechP.Angle
	s["cAngle"] = c.mechP.Cannon.Angle
	env.CreateAndInjectStruct("Mech", "mech", s)
}

func (c *Code) loadWorldObjectsIntoEnv(env *object.Environment) {
	a := make([]object.AbstractStruct, 0)
	for _, o := range c.worldP.objects {
		if o.getType() == TypeMissile {
			// for the time being load only static objects
			continue
		}
		f := make(map[string]interface{})
		// todo need support of strings in marslang
		//f["type"] = o.getType()
		f["x"] = o.getPos().x
		f["y"] = o.getPos().y
		f["angle"] = o.getAngle()
		a = append(a, object.AbstractStruct{Fields: f})
	}
	env.CreateAndInjectArrayOfStructs("Object", "objects", a)
}

func (c *Code) Run() {
	ticker := time.NewTicker(c.runSpeedMs * time.Millisecond)

	for range ticker.C {
		//log.Printf("Code run tick\n")
		c.listenTheWorld()
		if c.astProgram == nil || c.state != Running {
			// waiting for the ast or for the launch
			continue
		}
		env := object.NewEnvironment()
		c.loadMechVarsIntoEnv(env)
		c.loadWorldObjectsIntoEnv(env)

		_, err := interpereter.Exec(c.astProgram, env)
		if err != nil {
			c.state = Stopped
			c.errorCh <- &Error{
				ErrorType: Runtime,
				Message:   err.Error(),
			}
			log.Printf("Runtime error: %s", err.Error())
			continue
		}

		c.outputCh <- newMechOutputVarsFromEnv(env)

		env.Print()
	}
}

func (c *Code) listenTheWorld() {
	select {
	case c.state = <-c.flowCh:
		if c.state == Stopped {
			c.outputCh <- &MechOutputVars{
				MThrottle: 0,
				RThrottle: 0,
			}
		}
	case c.astProgram = <-c.codeSaveCh:
		log.Println("Code saved")
	default:
		// noop
	}
}

func (c *Code) saveCode(sourceCode string) {
	l := lexer.New(sourceCode)
	p, err := parser.New(l)
	if err != nil {
		c.errorCh <- &Error{
			ErrorType: Lexing,
			Message:   err.Error(),
		}
		log.Printf("Lexing error: %s", err.Error())
		return
	}
	astProgram, err := p.Parse()
	if err != nil {
		c.errorCh <- &Error{
			ErrorType: Parsing,
			Message:   err.Error(),
		}
		log.Printf("Parsing error: %s", err.Error())
		return
	}
	log.Println("Code parsed")
	c.codeSaveCh <- astProgram
}

func (c *Code) operateState(cmd server.ProgramFlowType) {
	switch cmd {
	case server.StartProgram:
		c.flowCh <- Running
	case server.StopProgram:
		c.flowCh <- Stopped
	}
}
