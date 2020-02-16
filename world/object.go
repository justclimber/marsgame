package world

import (
	"aakimov/marsgame/changelog"
	"aakimov/marsgame/physics"
	"math"
)

const (
	TypePlayer    = "player"
	TypeMissile   = "missile"
	TypeRock      = "rock"
	TypeEnemyMech = "enemy_mech"
	TypeXelon     = "xelon"
)

type IObject interface {
	run(world *World) *changelog.ChangeByObject
	setId(id int)
	getId() int
	getType() string
	getPos() physics.Point
	getObj() physics.Obj
	getAngle() float64
	getMoveDone() *physics.Vector
	isCollideWith(o1 IObject) bool
	getCollisionRadius() int
}

type Object struct {
	physics.Obj
}

func NewObject(id int, typeObj string, p physics.Point, colRadius int, angle, speed, aspeed, weight float64) Object {
	return Object{Obj: physics.Obj{
		Id:              id,
		Type:            typeObj,
		Pos:             p,
		CollisionRadius: colRadius,
		Angle:           angle,
		Speed:           speed,
		AngleSpeed:      aspeed,
		Weight:          weight,
		Velocity:        &physics.Vector{},
		Direction:       physics.MakeNormalVectorByAngle(angle),
	}}
}

func (o *Object) run(world *World) *changelog.ChangeByObject {
	if o.Speed == 0 {
		o.MoveDone = nil
		return nil
	}

	ch := &changelog.ChangeByObject{
		ObjType: o.Type,
		ObjId:   o.Id,
	}

	o.Lock()
	defer o.Unlock()

	if o.AngleSpeed != 0 {
		o.Angle += o.AngleSpeed
		if o.Angle > 2*math.Pi {
			o.Angle = o.Angle - 2*math.Pi
		} else if o.Angle < 0 {
			o.Angle = 2*math.Pi + o.Angle
		}
	}

	o.MoveDone = o.Pos.MoveForward(o.Angle, o.Speed)
	if o.Pos.CheckIfOutOfBounds(0, 0, float64(world.width), float64(world.height)) {
		ch.Delete = true
		return ch
	}
	newPos := o.Pos
	newSpeed := o.Speed
	newAngle := o.Angle
	ch.Pos = &newPos
	ch.Length = &newSpeed
	ch.Angle = &newAngle
	return ch
}

func (o *Object) getId() int                   { return o.Id }
func (o *Object) setId(id int)                 { o.Id = id }
func (o *Object) getObj() physics.Obj          { return o.Obj }
func (o *Object) getPos() physics.Point        { return o.Pos }
func (o *Object) getAngle() float64            { return o.Angle }
func (o *Object) getType() string              { return o.Type }
func (o *Object) getMoveDone() *physics.Vector { return o.MoveDone }
func (o *Object) getCollisionRadius() int      { return o.CollisionRadius }
func (o *Object) isCollideWith(o1 IObject) bool {
	return o.IsCollideWith(o1.getObj())
}
