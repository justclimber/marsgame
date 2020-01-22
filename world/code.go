package world

import (
	"aakimov/marsgame/server"
	"aakimov/marslang/ast"
	"aakimov/marslang/interpereter"
	"aakimov/marslang/lexer"
	"aakimov/marslang/object"
	"aakimov/marslang/parser"

	"encoding/json"
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
	s["x"] = c.mechP.Pos.X
	s["y"] = c.mechP.Pos.Y
	s["angle"] = c.mechP.Angle
	s["cAngle"] = c.mechP.Cannon.Angle
	env.CreateAndInjectStruct("Mech", "mech", s)
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

		c.mu.Lock()
		astProgram := c.astProgram
		c.mu.Unlock()

		_, err := interpereter.Exec(astProgram, env)
		if err != nil {
			c.state = Stopped
			c.errorCh <- &Error{
				ErrorType: Runtime,
				Message:   err.Error(),
			}
			log.Printf("Runtime error: %s", err.Error())
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
	default:
		// noop
	}
}

func (c *Code) saveAst(ast *ast.StatementsBlock) {
	// todo переделать на каналы
	c.mu.Lock()
	c.astProgram = ast
	c.mu.Unlock()
}

func (c *Code) operateState(cmd server.ProgramFlowType) {
	switch cmd {
	case server.StartProgram:
		c.flowCh <- Running
	case server.StopProgram:
		c.flowCh <- Stopped
	}
}

func ParseSourceCode(sourceCode string) (*ast.StatementsBlock, []byte) {
	// todo передалать на метод структуры, чтобы сразу парсил код и сохранял
	log.Printf("Program to parse: %s\n", sourceCode)

	l := lexer.New(sourceCode)
	p, err := parser.New(l)
	if err != nil {
		return nil, respondWithError(err.Error(), "Lexing")
	}
	astProgram, err := p.Parse()
	if err != nil {
		return nil, respondWithError(err.Error(), "Parsing")
	}

	return astProgram, nil
}

func respondWithError(msg, prefix string) []byte {
	log.Printf("%s error: %s\n", prefix, msg)
	return errorToJson(msg)
}

func errorToJson(msg string) []byte {
	errJson := make(map[string]string)
	errJson["error"] = msg
	errBytes, _ := json.Marshal(errJson)

	return errBytes
}
