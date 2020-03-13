package world

import (
	"aakimov/marsgame/physics"
	"aakimov/marsgame/wal"
	"time"
)

const (
	TypePlayer    = "player"
	TypeMissile   = "missile"
	TypeRock      = "rock"
	TypeEnemyMech = "enemy_mech"
	TypeXelon     = "xelon"
	TypeSpore     = "spore"
)

type IObject interface {
	run(world *World, timeDelta time.Duration, timeId int64)
	setId(id uint32)
	getId() uint32
	getType() string
	getPos() physics.Point
	getObj() physics.Obj
	getAngle() float64
	isCollideWith(o1 IObject) bool
	getCollisionRadius() int
	setObjectObserver(om *wal.ObjectObserver)
}

type Object struct {
	physics.Obj
	wal *wal.ObjectObserver
}

type Missile struct {
	Object
}

func NewObject(id uint32, typeObj string, p physics.Point, colRadius int, angle, speed, aspeed, weight float64) Object {
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

func (o *Object) getId() uint32           { return o.Id }
func (o *Object) setId(id uint32)         { o.Id = id }
func (o *Object) getObj() physics.Obj     { return o.Obj }
func (o *Object) getPos() physics.Point   { return o.Pos }
func (o *Object) getAngle() float64       { return o.Angle }
func (o *Object) getType() string         { return o.Type }
func (o *Object) getCollisionRadius() int { return o.CollisionRadius }
func (o *Object) isCollideWith(o1 IObject) bool {
	return o.IsCollideWith(o1.getObj())
}
func (o *Object) setObjectObserver(om *wal.ObjectObserver) {
	o.wal = om
}
