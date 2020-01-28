package world

import (
	"aakimov/marsgame/server"
	"aakimov/marslang/ast"
	"aakimov/marslang/interpereter"
	"aakimov/marslang/lexer"
	"aakimov/marslang/object"
	"aakimov/marslang/parser"
	"fmt"
	"log"
	"math"
	"sort"
	"sync"
	"time"
)

type ProgramState int

const (
	Stopped ProgramState = iota
	Running
)

type Code struct {
	id          string
	state       ProgramState
	mu          sync.Mutex
	astProgram  *ast.StatementsBlock
	worldP      *World
	mechP       *Mech
	outputCh    chan *MechOutputVars
	io4ClientCh chan *IO4Client
	codeSaveCh  chan *ast.StatementsBlock
	flowCh      chan ProgramState
	errorCh     chan *Error
	runSpeedMs  time.Duration
	energy      int
}

func NewCode(id string, world *World, mech *Mech, runSpeedMs time.Duration) *Code {
	return &Code{
		id:          "main",
		worldP:      world,
		mechP:       mech,
		outputCh:    make(chan *MechOutputVars),
		codeSaveCh:  make(chan *ast.StatementsBlock),
		io4ClientCh: make(chan *IO4Client),
		flowCh:      make(chan ProgramState),
		errorCh:     make(chan *Error),
		runSpeedMs:  runSpeedMs,
		energy:      10000,
	}
}

type IO4Client struct {
	Input  []string
	Output []string
	Cost   int
	Energy int
}

type MechOutputVars struct {
	MThrottle  float64
	RThrottle  float64
	CRThrottle float64
	Shoot      float64
}

func (m *MechOutputVars) toStrings() []string {
	result := make([]string, 0)
	result = append(result, fmt.Sprintf("mThr = %.2f", m.MThrottle))
	result = append(result, fmt.Sprintf("mrThr = %.2f", m.RThrottle))
	result = append(result, fmt.Sprintf("crThr = %.2f", m.CRThrottle))
	result = append(result, fmt.Sprintf("shoot = %.2f", m.Shoot))
	return result
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

func (c *Code) loadMiscIntoEnv(env *object.Environment) {
	env.Set("PI", &object.Float{Value: math.Pi})
}

func (c *Code) Run() {
	ticker := time.NewTicker(c.runSpeedMs * time.Millisecond)

	executor := interpereter.NewExecAstVisitor()
	executor.SetExecCallback(c.consumeEnergy)
	lastEnergy := c.energy
	env := object.NewEnvironment()
	for range ticker.C {
		//log.Printf("Code run tick\n")
		c.listenTheWorld()
		if c.astProgram == nil || c.state != Running {
			// waiting for the ast or for the launch
			continue
		}

		c.loadMechVarsIntoEnv(env)
		c.loadWorldObjectsIntoEnv(env)
		c.loadMiscIntoEnv(env)

		err := executor.ExecAst(c.astProgram, env)
		if err != nil {
			c.state = Stopped
			c.errorCh <- &Error{
				ErrorType: Runtime,
				Message:   err.Error(),
			}
			log.Printf("Runtime error: %s", err.Error())
			continue
		}
		cost := lastEnergy - c.energy
		lastEnergy = c.energy

		c.io4ClientCh <- c.makeIO4Client(env, cost, c.energy)
		c.outputCh <- newMechOutputVarsFromEnv(env)
	}
}

var energyByOperationMap = map[interpereter.OperationType]int{
	interpereter.Assignment:      20,
	interpereter.Return:          15,
	interpereter.IfStmt:          15,
	interpereter.Switch:          15,
	interpereter.Unary:           4,
	interpereter.BinExpr:         20,
	interpereter.Struct:          25,
	interpereter.StructFieldCall: 8,
	interpereter.NumInt:          3,
	interpereter.NumFloat:        4,
	interpereter.Boolean:         2,
	interpereter.Array:           6,
	interpereter.ArrayIndex:      4,
	interpereter.Identifier:      6,
	interpereter.Function:        15,
	interpereter.FunctionCall:    10,
}

func (c *Code) consumeEnergy(operation interpereter.OperationType) {
	energyCost, ok := energyByOperationMap[operation]
	if !ok {
		log.Fatalf("Unknown operation for energy calculation: %v", operation)
	}
	c.energy -= energyCost
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

func (c *Code) makeIO4Client(env *object.Environment, cost, energy int) *IO4Client {
	inputKeys := map[string]bool{"mech": true, "objects": true, "PI": true}
	input := make([]string, 0)
	output := make([]string, 0)
	for k, v := range env.Store() {
		vStr := fmt.Sprintf("%s: %s\n", k, v.Inspect())
		if _, ok := inputKeys[k]; ok {
			input = append(input, vStr)
		} else {
			output = append(output, vStr)
		}
	}
	sort.Strings(input)
	sort.Strings(output)
	return &IO4Client{
		Input:  input,
		Output: output,
		Cost:   cost,
		Energy: energy,
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
