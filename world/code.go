package world

import (
	"aakimov/marslang/ast"
	"aakimov/marslang/interpereter"
	"aakimov/marslang/object"
	"log"
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
	var objTypeToIntMap = map[string]int{
		TypePlayer:    0,
		TypeEnemyMech: 1,
		TypeRock:      2,
		TypeXelon:     3,
		TypeMissile:   4,
	}
	a := make([]object.AbstractStruct, 0)
	for _, o := range w.objects {
		if o.getType() == TypeMissile {
			// for the time being load only static objects
			continue
		}
		f := make(map[string]interface{})
		f["type"] = objTypeToIntMap[o.getType()]
		f["x"] = o.getPos().X
		f["y"] = o.getPos().Y
		f["angle"] = o.getAngle()
		a = append(a, object.AbstractStruct{Fields: f})
	}
	if len(a) == 0 {
		f := make(map[string]interface{})
		f["type"] = 0
		f["x"] = 0.
		f["y"] = 0.
		f["angle"] = 0.
		a = append(a, object.AbstractStruct{Fields: f})
		env.CreateAndInjectEmptyArrayOfStructs("Object", "objects", a)
		log.Println("EpmtyObjectsLoaded")
		return
	}
	env.CreateAndInjectArrayOfStructs("Object", "objects", a)
}

func (c *Code) loadMiscIntoEnv(env *object.Environment) {
	env.Set("PI", &object.Float{Value: math.Pi})
}

const (
	bDistance      string = "distance"
	bAngle         string = "angle"
	bNearest       string = "nearest"
	bNearestByType string = "nearestByType"
)

func (c *Code) SetupMarsGameBuiltinFunctions(executor *interpereter.ExecAstVisitor) {
	builtins := make(map[string]*object.Builtin)
	builtins[bDistance] = &object.Builtin{
		Name:       bDistance,
		ReturnType: object.FloatObj,
		Fn: func(env *object.Environment, args ...object.Object) (object.Object, error) {
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
	}
	builtins[bAngle] = &object.Builtin{
		Name:       bAngle,
		ReturnType: object.FloatObj,
		Fn: func(env *object.Environment, args ...object.Object) (object.Object, error) {
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
	}
	builtins[bNearest] = &object.Builtin{
		Name:       bNearest,
		ReturnType: "Object",
		Fn: func(env *object.Environment, args ...object.Object) (object.Object, error) {
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
			def, ok := env.GetStructDefinition("Object")
			if !ok {
				return nil, interpereter.BuiltinFuncError("Why no Object struct defined??")
			}
			if arrayOfStruct.Empty {
				return object.NewEmptyStruct(def), nil
			}
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
			return arrayOfStruct.Elements[minIndex], nil
		},
	}
	builtins[bNearestByType] = &object.Builtin{
		Name:       bNearestByType,
		ReturnType: "Object",
		Fn: func(env *object.Environment, args ...object.Object) (object.Object, error) {
			if len(args) != 3 {
				return nil, interpereter.BuiltinFuncError("wrong number of arguments. got=%d, want 3", len(args))
			}
			if err := interpereter.CheckArgType("Mech", args[0]); err != nil {
				return nil, err
			}
			if err := interpereter.CheckArgType("Object[]", args[1]); err != nil {
				return nil, err
			}
			if err := interpereter.CheckArgType("int", args[2]); err != nil {
				return nil, err
			}
			mech := args[0].(*object.Struct)
			arrayOfStruct, _ := args[1].(*object.Array)
			def, ok := env.GetStructDefinition("Object")
			if !ok {
				return nil, interpereter.BuiltinFuncError("Why no Object struct defined??")
			}
			if arrayOfStruct.Empty {
				return object.NewEmptyStruct(def), nil
			}
			objType := args[2].(*object.Integer).Value
			minDist := 99999999999.
			minIndex := -1
			for i := 0; i < len(arrayOfStruct.Elements); i++ {
				obj, _ := arrayOfStruct.Elements[i].(*object.Struct)
				if obj.Fields["type"].(*object.Integer).Value != objType {
					continue
				}
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
				return object.NewEmptyStruct(def), nil
			}
			return arrayOfStruct.Elements[minIndex], nil
		},
	}
	executor.AddBuiltinFunctions(builtins)
}
