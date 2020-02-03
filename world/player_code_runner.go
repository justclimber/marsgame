package world

import (
	"aakimov/marsgame/server"
	"aakimov/marslang/interpereter"
	"aakimov/marslang/object"

	"fmt"
	"log"
	"sort"
	"time"
)

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

func (p *Player) runProgram() {
	code := p.mainProgram

	executor := interpereter.NewExecAstVisitor()
	executor.SetExecCallback(p.consumeEnergy)
	code.SetupMarsGameBuiltinFunctions(executor)

	go p.mech.generator.start()

	ticker := time.NewTicker(p.runSpeedMs * time.Millisecond)
	for range ticker.C {
		//log.Printf("Code runProgram tick\n")
		p.listenBeforeCodeExec()
		if code.astProgram == nil || code.state != Running {
			// waiting for the ast or for the launch
			continue
		}
		env := object.NewEnvironment()
		code.bootstrap(p, env)

		err := executor.ExecAst(code.astProgram, env)
		if err != nil {
			code.state = Stopped
			p.errorCh <- &Error{
				ErrorType: Runtime,
				Message:   err.Error(),
			}
			log.Printf("Runtime error: %s", err.Error())
			continue
		}

		p.io4ClientCh <- p.makeIO4Client(env)
		p.outputCh <- newMechOutputVarsFromEnv(env)

		p.mainProgram.codeExecCost = 0
	}
}

func (p *Player) consumeEnergy(operation interpereter.Operation) {
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
	var energyByBuiltinFunctionMap = map[string]int{
		bDistance: 60,
		bAngle:    50,
		bNearest:  200,
		"print":   10,
	}
	var ok bool
	var energyCost int

	if operation.Type == interpereter.Builtin {
		energyCost, ok = energyByBuiltinFunctionMap[operation.FuncName]
	} else {
		energyCost, ok = energyByOperationMap[operation.Type]
	}

	if !ok {
		log.Fatalf("Unknown operation for energy calculation: %+v\n", operation)
	}
	p.mech.generator.consumeWithThrottling(energyCost)
	p.mainProgram.codeExecCost += energyCost
}

func (p *Player) listenBeforeCodeExec() {
	select {
	case p.mainProgram.state = <-p.flowCh:
		if p.mainProgram.state == Stopped {
			p.outputCh <- &MechOutputVars{
				MThrottle: 0,
				RThrottle: 0,
			}
		}
	case p.mainProgram.astProgram = <-p.codeSaveCh:
		log.Println("Code saved")
	default:
	}
}

func (p *Player) makeIO4Client(env *object.Environment) *IO4Client {
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
		Cost:   p.mainProgram.codeExecCost,
		Energy: p.mech.generator.geValue(),
	}
}

func (p *Player) operateState(cmd server.ProgramFlowType) {
	switch cmd {
	case server.StartProgram:
		p.flowCh <- Running
	case server.StopProgram:
		p.flowCh <- Stopped
	}
}
