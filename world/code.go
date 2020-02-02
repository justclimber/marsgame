package world

import (
	"aakimov/marslang/ast"
	"aakimov/marslang/interpereter"
	"aakimov/marslang/object"
	"math"
	"sync"
	"time"
)

type ProgramState int

const (
	Stopped ProgramState = iota
	Running
)

type Code struct {
	id           string
	state        ProgramState
	mu           sync.Mutex
	astProgram   *ast.StatementsBlock
	codeExecCost int
}

func NewCode(id string, world *World, mech *Mech, runSpeedMs time.Duration) *Code {
	return &Code{
		id: "main",
	}
}

func (c *Code) bootstrap(p *Player, env *object.Environment) {
	c.loadMechVarsIntoEnv(p.mech, env)
	c.loadWorldObjectsIntoEnv(p.world, env)
	c.loadMiscIntoEnv(env)
}

func (c *Code) loadMechVarsIntoEnv(m *Mech, env *object.Environment) {
	s := make(map[string]interface{})
	s["x"] = m.Pos.X
	s["y"] = m.Pos.Y
	s["angle"] = m.Angle
	s["cAngle"] = m.cannon.angle
	env.CreateAndInjectStruct("Mech", "mech", s)
}

func (c *Code) loadWorldObjectsIntoEnv(w *World, env *object.Environment) {
	a := make([]object.AbstractStruct, 0)
	for _, o := range w.objects {
		if o.getType() == TypeMissile {
			// for the time being load only static objects
			continue
		}
		f := make(map[string]interface{})
		// todo need support of strings in marslang
		//f["type"] = o.getType()
		f["x"] = o.getPos().X
		f["y"] = o.getPos().Y
		f["angle"] = o.getAngle()
		a = append(a, object.AbstractStruct{Fields: f})
	}
	env.CreateAndInjectArrayOfStructs("Object", "objects", a)
}

func (c *Code) loadMiscIntoEnv(env *object.Environment) {
	env.Set("PI", &object.Float{Value: math.Pi})
}

func (c *Code) SetupMarsGameBuiltinFunctions(executor *interpereter.ExecAstVisitor) {
	builtins := make(map[string]*object.Builtin)
	builtins["distance"] = &object.Builtin{
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) != 4 {
				return nil, interpereter.BuiltinFuncError("wrong number of arguments. got=%d, want 2", len(args))
			}
			if err := interpereter.CheckArgsType(object.FloatObj, args); err != nil {
				return nil, err
			}
			x1 := args[0].(*object.Float).Value
			y1 := args[1].(*object.Float).Value
			x2 := args[2].(*object.Float).Value
			y2 := args[3].(*object.Float).Value
			return &object.Float{Value: distance(x1, y1, x2, y2)}, nil
		},
		ReturnType: object.FloatObj,
	}
	builtins["angle"] = &object.Builtin{
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) != 4 {
				return nil, interpereter.BuiltinFuncError("wrong number of arguments. got=%d, want 2", len(args))
			}
			if err := interpereter.CheckArgsType(object.FloatObj, args); err != nil {
				return nil, err
			}
			x1 := args[0].(*object.Float).Value
			y1 := args[1].(*object.Float).Value
			x2 := args[2].(*object.Float).Value
			y2 := args[3].(*object.Float).Value
			return &object.Float{Value: angle(x1, y1, x2, y2)}, nil
		},
		ReturnType: object.FloatObj,
	}
	builtins["nearest"] = &object.Builtin{
		Fn: func(args ...object.Object) (object.Object, error) {
			if len(args) != 2 {
				return nil, interpereter.BuiltinFuncError("wrong number of arguments. got=%d, want 2", len(args))
			}
			if err := interpereter.CheckArgType("Mech", args[0]); err != nil {
				return nil, err
			}
			if err := interpereter.CheckArgType("Object[]", args[1]); err != nil {
				return nil, err
			}
			mech := args[0].(*object.Struct)
			arrayOfStruct, _ := args[1].(*object.Array)
			minDist := 99999999999.
			minIndex := -1
			for i := 0; i < len(arrayOfStruct.Elements); i++ {
				obj, _ := arrayOfStruct.Elements[i].(*object.Struct)
				objX := obj.Fields["x"].(*object.Float).Value
				objY := obj.Fields["y"].(*object.Float).Value
				mechX := mech.Fields["x"].(*object.Float).Value
				mechY := mech.Fields["y"].(*object.Float).Value
				dist := distance(mechX, mechY, objX, objY)
				if dist < minDist {
					minDist = dist
					minIndex = i
				}
			}
			if minIndex == -1 {
				return nil, interpereter.BuiltinFuncError("nearest on empty array")
			}
			return arrayOfStruct.Elements[minIndex], nil
		},
		ReturnType: "Object",
	}
	executor.AddBuiltinFunctions(builtins)
}
