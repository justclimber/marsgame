package world

import (
	"strconv"
	"sync"
)

const (
	TypeObject  = "object"
	TypePlayer  = "player"
	TypeMissile = "missile"
	TypeRock    = "rock"
)

type IObject interface {
	run(world *World) *ChangeByObject
	setId(id int)
	getId() int
	getType() string
	getPos() Point
	getAngle() float64
	getMoveDone() *Vector
	isCollideWith(o1 IObject) bool
	getCollisionRadius() int
}

type Object struct {
	Id              int
	Type            string
	mu              sync.Mutex
	Pos             Point
	CollisionRadius int
	Angle           float64
	Speed           float64
	AngleSpeed      float64
	MoveDone        *Vector
}

func (o *Object) run(world *World) *ChangeByObject {
	if o.Speed == 0 {
		o.MoveDone = nil
		return nil
	}

	ch := &ChangeByObject{
		ObjType: o.Type,
		ObjId:   strconv.Itoa(o.Id),
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	o.MoveDone = o.Pos.MoveForward(o.Angle, o.Speed)
	if o.Pos.checkIfOutOfBounds(0, 0, float64(world.width), float64(world.height)) {
		ch.Delete = true
		return ch
	}
	newPos := o.Pos
	newSpeed := o.Speed
	ch.Pos = &newPos
	ch.length = &newSpeed
	ch.Angle = &o.Angle
	return ch
}

func (o *Object) getId() int              { return o.Id }
func (o *Object) setId(id int)            { o.Id = id }
func (o *Object) getPos() Point           { return o.Pos }
func (o *Object) getAngle() float64       { return o.Angle }
func (o *Object) getType() string         { return o.Type }
func (o *Object) getMoveDone() *Vector    { return o.MoveDone }
func (o *Object) getCollisionRadius() int { return o.CollisionRadius }
