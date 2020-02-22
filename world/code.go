package world

import (
	"aakimov/marsgame/physics"
	"aakimov/marslang/ast"
	"aakimov/marslang/interpereter"
	"aakimov/marslang/object"
	"math"
	"sync"
)

type ProgramState int

const (
	Stopped ProgramState = iota
	Running
	Terminate
)

type Code struct {
	id               string
	state            ProgramState
	mu               sync.Mutex
	astProgram       *ast.StatementsBlock
	codeExecCost     int
	objTargetsByType ObjTargetsByType
}

type ObjTargets []*object.Struct
type ObjTargetsByType struct {
	targets map[int]ObjTargets
}

func (o *ObjTargetsByType) Add(targetType int, obj *object.Struct) {
	_, exist := o.targets[targetType]
	if !exist {
		o.targets[targetType] = make([]*object.Struct, 0)
	}
	o.targets[targetType] = append(o.targets[targetType], obj)
}

func (o *ObjTargetsByType) GetFirst(targetType int, def *object.StructDefinition) *object.Struct {
	targetsTyped, exist := o.targets[targetType]
	if !exist || len(targetsTyped) == 0 {
		return object.NewEmptyStruct(def)
	}
	return targetsTyped[0]
}

func (o *ObjTargetsByType) actualize(check map[int]bool) {
	for k, v := range o.targets {
		newTargets := make([]*object.Struct, 0)
		for _, vv := range v {
			id := vv.Fields["id"].(*object.Integer).Value
			if _, exist := check[int(id)]; exist {
				newTargets = append(newTargets, vv)
			}
		}
		if len(newTargets) == 0 {
			delete(o.targets, k)
		} else if len(newTargets) != len(o.targets[k]) {
			o.targets[k] = newTargets
		}
	}
}

func NewCode(id string) *Code {
	return &Code{
		id:               id,
		objTargetsByType: ObjTargetsByType{targets: make(map[int]ObjTargets)},
	}
}

func (c *Code) bootstrap(p *Player, structs map[string]*object.StructDefinition, env *object.Environment) {
	c.copyStructDefinitionsToEnv(structs, env)
	c.loadMechVarsIntoEnv(p.mech, structs, env)
	c.loadWorldObjectsIntoEnv(p.world, structs, env)
	c.loadMiscIntoEnv(env)
	c.loadCommandsVarIntoEnv(structs, env)
}

func (c *Code) getStructDefinitions() map[string]*object.StructDefinition {
	return map[string]*object.StructDefinition{
		"Mech": {"Mech", map[string]string{
			"x":      object.FloatObj,
			"y":      object.FloatObj,
			"angle":  object.FloatObj,
			"cAngle": object.FloatObj,
		}},
		"Object": {"Object", map[string]string{
			"id":    object.IntegerObj,
			"type":  object.IntegerObj,
			"x":     object.FloatObj,
			"y":     object.FloatObj,
			"angle": object.FloatObj,
		}},
		"Cannon": {"Cannon", map[string]string{
			"rotate": object.FloatObj,
			"shoot":  object.FloatObj,
		}},
		"Commands": {"Commands", map[string]string{
			"move":   object.FloatObj,
			"rotate": object.FloatObj,
			"cannon": "Cannon",
		}},
	}
}

func (c *Code) copyStructDefinitionsToEnv(structs map[string]*object.StructDefinition, env *object.Environment) {
	for _, v := range structs {
		_ = env.RegisterStructDefinition(v)
	}
}

func (c *Code) loadCommandsVarIntoEnv(structs map[string]*object.StructDefinition, env *object.Environment) {
	commands := env.LoadVarsInStruct(structs["Commands"], map[string]interface{}{
		"move":   0.,
		"rotate": 0.,
	})
	commands.Fields["cannon"] = env.LoadVarsInStruct(structs["Cannon"], map[string]interface{}{
		"rotate": 0.,
		"shoot":  0.,
	})
	env.Set("commands", commands)
}

func (c *Code) loadMechVarsIntoEnv(m *Mech, structs map[string]*object.StructDefinition, env *object.Environment) {
	s := env.LoadVarsInStruct(structs["Mech"], map[string]interface{}{
		"x":      m.Pos.X,
		"y":      m.Pos.Y,
		"angle":  m.Angle,
		"cAngle": m.cannon.angle,
	})
	env.Set("mech", s)
}

func (c *Code) loadWorldObjectsIntoEnv(w *World, structs map[string]*object.StructDefinition, env *object.Environment) {
	var objTypeToIntMap = map[string]int{
		TypePlayer:    0,
		TypeEnemyMech: 1,
		TypeRock:      2,
		TypeXelon:     3,
		TypeMissile:   4,
	}
	elements := make([]object.Object, 0)
	check := make(map[int]bool)
	for _, o := range w.objects {
		if o.getType() == TypeMissile {
			// for the time being load only static objects
			continue
		}
		check[o.getId()] = true
		elements = append(elements, env.LoadVarsInStruct(structs["Object"], map[string]interface{}{
			"id":    o.getId(),
			"type":  objTypeToIntMap[o.getType()],
			"x":     o.getPos().X,
			"y":     o.getPos().Y,
			"angle": o.getAngle(),
		}))
	}

	c.objTargetsByType.actualize(check)

	resultArray := &object.Array{ElementsType: "Object"}
	if len(elements) == 0 {
		resultArray.Empty = true
		resultArray.Elements = make([]object.Object, 0)
	} else {
		resultArray.Elements = elements
	}
	env.Set("objects", resultArray)
}

func (c *Code) loadMiscIntoEnv(env *object.Environment) {
	env.Set("PI", &object.Float{Value: math.Pi})
}

const (
	bDistance          string = "distance"
	bAngle             string = "angle"
	bAngleToRotate     string = "angleToRotate"
	bNearest           string = "nearest"
	bNearestByType     string = "nearestByType"
	bAddTarget         string = "addTarget"
	bGetFirstTarget    string = "getFirstTarget"
	bRemoveFirstTarget string = "removeFirstTarget"
	bKeepBounds        string = "keepBounds"
)

func (c *Code) SetupMarsGameBuiltinFunctions(
	executor *interpereter.ExecAstVisitor,
	structDefs map[string]*object.StructDefinition,
) {
	builtins := make(map[string]*object.Builtin)
	builtins[bDistance] = &object.Builtin{
		Name:       bDistance,
		ArgTypes:   object.ArgTypes{object.FloatObj, object.FloatObj, object.FloatObj, object.FloatObj},
		ReturnType: object.FloatObj,
		Fn: func(env *object.Environment, args []object.Object) (object.Object, error) {
			x1 := args[0].(*object.Float).Value
			y1 := args[1].(*object.Float).Value
			x2 := args[2].(*object.Float).Value
			y2 := args[3].(*object.Float).Value
			return &object.Float{Value: physics.Distance(x1, y1, x2, y2)}, nil
		},
	}
	builtins[bAngle] = &object.Builtin{
		Name:       bAngle,
		ArgTypes:   object.ArgTypes{object.FloatObj, object.FloatObj, object.FloatObj, object.FloatObj},
		ReturnType: object.FloatObj,
		Fn: func(env *object.Environment, args []object.Object) (object.Object, error) {
			x1 := args[0].(*object.Float).Value
			y1 := args[1].(*object.Float).Value
			x2 := args[2].(*object.Float).Value
			y2 := args[3].(*object.Float).Value
			return &object.Float{Value: physics.Angle(x1, y1, x2, y2)}, nil
		},
	}
	builtins[bAngleToRotate] = &object.Builtin{
		Name:       bAngleToRotate,
		ArgTypes:   object.ArgTypes{object.FloatObj, object.FloatObj, object.FloatObj, object.FloatObj, object.FloatObj},
		ReturnType: object.FloatObj,
		Fn: func(env *object.Environment, args []object.Object) (object.Object, error) {
			angleFrom := args[0].(*object.Float).Value
			x1 := args[1].(*object.Float).Value
			y1 := args[2].(*object.Float).Value
			x2 := args[3].(*object.Float).Value
			y2 := args[4].(*object.Float).Value
			if angleFrom > 2*math.Pi {
				angleFrom -= 2 * math.Pi
			}
			angleTo := physics.Angle(x1, y1, x2, y2) - angleFrom
			if angleTo < -math.Pi {
				angleTo = 2.*math.Pi + angleTo
			} else if angleTo > math.Pi {
				angleTo = angleTo - 2.*math.Pi
			}
			return &object.Float{Value: angleTo}, nil
		},
	}
	builtins[bNearest] = &object.Builtin{
		Name:       bNearest,
		ArgTypes:   object.ArgTypes{"Mech", "Object[]"},
		ReturnType: "Object",
		Fn: func(env *object.Environment, args []object.Object) (object.Object, error) {
			mech := args[0].(*object.Struct)
			arrayOfStruct, _ := args[1].(*object.Array)
			if arrayOfStruct.Empty {
				return object.NewEmptyStruct(structDefs["Object"]), nil
			}
			minDist := 99999999999.
			minIndex := -1
			for i := 0; i < len(arrayOfStruct.Elements); i++ {
				obj, _ := arrayOfStruct.Elements[i].(*object.Struct)
				objX := obj.Fields["x"].(*object.Float).Value
				objY := obj.Fields["y"].(*object.Float).Value
				mechX := mech.Fields["x"].(*object.Float).Value
				mechY := mech.Fields["y"].(*object.Float).Value
				dist := physics.Distance(mechX, mechY, objX, objY)
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
		ArgTypes:   object.ArgTypes{"Mech", "Object[]", object.IntegerObj},
		ReturnType: "Object",
		Fn: func(env *object.Environment, args []object.Object) (object.Object, error) {
			mech := args[0].(*object.Struct)
			arrayOfStruct, _ := args[1].(*object.Array)
			if arrayOfStruct.Empty {
				return object.NewEmptyStruct(structDefs["Object"]), nil
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
				dist := physics.Distance(mechX, mechY, objX, objY)
				if dist < minDist {
					minDist = dist
					minIndex = i
				}
			}
			if minIndex == -1 {
				return object.NewEmptyStruct(structDefs["Object"]), nil
			}
			return arrayOfStruct.Elements[minIndex], nil
		},
	}
	builtins[bAddTarget] = &object.Builtin{
		Name:       bAddTarget,
		ArgTypes:   object.ArgTypes{"Object", object.IntegerObj},
		ReturnType: object.VoidObj,
		Fn: func(env *object.Environment, args []object.Object) (object.Object, error) {
			objStruct, _ := args[0].(*object.Struct)
			targetType := int(args[1].(*object.Integer).Value)
			c.objTargetsByType.Add(targetType, objStruct)

			return &object.Void{}, nil
		},
	}
	builtins[bGetFirstTarget] = &object.Builtin{
		Name:       bGetFirstTarget,
		ArgTypes:   object.ArgTypes{object.IntegerObj},
		ReturnType: "Object",
		Fn: func(env *object.Environment, args []object.Object) (object.Object, error) {
			targetType := int(args[0].(*object.Integer).Value)

			return c.objTargetsByType.GetFirst(targetType, structDefs["Object"]), nil
		},
	}
	builtins[bKeepBounds] = &object.Builtin{
		Name:       bKeepBounds,
		ArgTypes:   object.ArgTypes{object.FloatObj, object.FloatObj},
		ReturnType: object.FloatObj,
		Fn: func(env *object.Environment, args []object.Object) (object.Object, error) {
			value := args[0].(*object.Float).Value
			bound := args[1].(*object.Float).Value
			if value > bound {
				value = bound
			} else if value < -bound {
				value = -bound
			}

			return &object.Float{Value: value}, nil
		},
	}
	executor.AddBuiltinFunctions(builtins)
}
