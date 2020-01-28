package world

import (
	"aakimov/marslang/ast"
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
	s["x"] = m.Pos.x
	s["y"] = m.Pos.y
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
