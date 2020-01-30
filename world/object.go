package world

import (
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
	Id   int
	Type string
	sync.Mutex
	Pos             Point
	CollisionRadius int
	Angle           float64
	Speed           float64
	AngleSpeed      float64
	Weight          float64
	MoveDone        *Vector
	Velocity        *Vector
	Direction       *Vector
}

func NewObject(id int, typeObj string, p Point, colRadius int, angle, speed, aspeed, weight float64) Object {
	return Object{
		Id:              id,
		Type:            typeObj,
		Pos:             p,
		CollisionRadius: colRadius,
		Angle:           angle,
		Speed:           speed,
		AngleSpeed:      aspeed,
		Weight:          weight,
		Velocity:        &Vector{},
		Direction:       makeNormalVectorByAngle(angle),
	}
}

func (o *Object) run(world *World) *ChangeByObject {
	if o.Speed == 0 {
		o.MoveDone = nil
		return nil
	}

	ch := &ChangeByObject{
		ObjType: o.Type,
		ObjId:   o.Id,
	}

	o.Lock()
	defer o.Unlock()

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
