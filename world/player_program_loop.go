package world

import (
	"github.com/justclimber/marsgame/server"
	"github.com/justclimber/marslang/interpereter"
	"github.com/justclimber/marslang/object"

	"fmt"
	"log"
	"sort"
	"time"
)

type IO4Client struct {
	Input    []string
	Output   []string
	Commands []string
	Cost     int
	Energy   int
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

func newMechOutputVarsFromEnv(env *object.Environment) *Commands {
	commands, _ := env.Get("commands")
	cmd := commands.(*object.Struct)
	cannon := cmd.Fields["cannon"].(*object.Struct)

	shootValue := cannon.Fields["shoot"].(*object.Float).Value
	shoot := &Shoot{}
	if shootValue != 0 {
		shoot.state = WillShoot
		shoot.delay = int(shootValue * 1000)
	}
	return &Commands{
		mech: &MechCommands{
			move:   cmd.Fields["move"].(*object.Float).Value,
			rotate: cmd.Fields["rotate"].(*object.Float).Value,
		},
		cannon: &CannonCommands{
			rotate: cannon.Fields["rotate"].(*object.Float).Value,
			shoot:  shoot,
		},
	}
}

func (p *Player) programLoop() {
	code := p.mainProgram

	executor := interpereter.NewExecAstVisitor()
	executor.SetExecCallback(p.consumeEnergy)

	structs := code.getStructDefinitions()
	enums := code.getEnumDefinitions()
	code.SetupMarsGameBuiltinFunctions(executor, structs)

	go p.mech.generator.start()

	ticker := time.NewTicker(p.runSpeedMs * time.Millisecond)
	for range ticker.C {
		//log.Printf("Code programLoop tick\n")
		select {
		case p.mainProgram.state = <-p.flowCh:
			switch p.mainProgram.state {
			case Stopped:
				p.codeCommandsCh <- NewEmptyCommands()
			case Terminate:
				p.codeCommandsCh <- NewEmptyCommands()
				p.mech.generator.terminate()
				p.terminateCh <- true
				ticker.Stop()
				return
			}
		case p.mainProgram.astProgram = <-p.codeSaveCh:
			log.Println("Code saved")
		default:
		}
		if code.astProgram == nil || code.state != Running {
			// waiting for the ast or for the launch
			continue
		}
		env := object.NewEnvironment()
		code.bootstrap(p, structs, enums, env)

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
		p.codeCommandsCh <- newMechOutputVarsFromEnv(env)

		p.mainProgram.codeExecCost = 0
	}
}

func (p *Player) consumeEnergy(operation interpereter.Operation) {
	var energyByOperationMap = map[interpereter.OperationType]int{
		interpereter.StructFieldAssignment: 20,
		interpereter.Assignment:            20,
		interpereter.Return:                15,
		interpereter.IfStmt:                15,
		interpereter.Switch:                15,
		interpereter.Unary:                 4,
		interpereter.BinExpr:               20,
		interpereter.Struct:                25,
		interpereter.StructFieldCall:       8,
		interpereter.NumInt:                3,
		interpereter.NumFloat:              4,
		interpereter.Boolean:               2,
		interpereter.Array:                 6,
		interpereter.ArrayIndex:            4,
		interpereter.Identifier:            6,
		interpereter.Function:              15,
		interpereter.FunctionCall:          10,
		interpereter.EnumElementCall:       4,
	}
	var energyByBuiltinFunctionMap = map[string]int{
		bDistance:                    60,
		bAngle:                       50,
		bAngleToRotate:               80,
		bNearest:                     300,
		bNearestByType:               350,
		bAddTarget:                   30,
		bGetFirstTarget:              30,
		bRemoveFirstTarget:           20,
		bKeepBounds:                  20,
		interpereter.BuiltinPrint:    10,
		interpereter.BuiltinEmpty:    10,
		interpereter.BuiltinAbsInt:   5,
		interpereter.BuiltinAbsFloat: 5,
		interpereter.BuiltinLength:   7,
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

func (p *Player) makeIO4Client(env *object.Environment) *IO4Client {
	inputKeys := map[string]bool{"mech": true, "objects": true, "PI": true}
	input := make([]string, 0)
	output := make([]string, 0)
	commands := make([]string, 4)
	for k, v := range env.Store() {
		vStr := fmt.Sprintf("%s: %s\n", k, v.Inspect())
		if _, ok := inputKeys[k]; ok {
			input = append(input, vStr)
		} else {
			if k == "commands" {
				commandsStruct, _ := v.(*object.Struct)
				commands[0] = fmt.Sprintf("move: %s", commandsStruct.Fields["move"].Inspect())
				commands[1] = fmt.Sprintf("rotate: %s", commandsStruct.Fields["rotate"].Inspect())
				cannonCommandsStruct, _ := commandsStruct.Fields["cannon"].(*object.Struct)
				commands[2] = fmt.Sprintf("cannon.rotate: %s", cannonCommandsStruct.Fields["rotate"].Inspect())
				commands[3] = fmt.Sprintf("cannon.shoot: %s", cannonCommandsStruct.Fields["shoot"].Inspect())
			} else {
				output = append(output, vStr)
			}
		}
	}
	sort.Strings(input)
	sort.Strings(output)
	return &IO4Client{
		Input:    input,
		Output:   output,
		Commands: commands,
		Cost:     p.mainProgram.codeExecCost,
		Energy:   p.mech.generator.geValue(),
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
