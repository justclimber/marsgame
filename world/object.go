package world

import (
	"aakimov/marsgame/physics"
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
	getPos() physics.Point
}

type Object struct {
	Id         int
	Type       string
	mu         sync.Mutex
	Pos        physics.Point
	Width      int
	Angle      float64
	Speed      float64
	AngleSpeed float64
}

func (o *Object) run(world *World) *ChangeByObject {
	if o.Speed == 0 {
		return nil
	}

	ch := &ChangeByObject{
		ObjType: o.Type,
		ObjId:   strconv.Itoa(o.Id),
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	o.Pos.MoveForward(o.Angle, o.Speed)
	if o.Pos.CheckIfOutOfBounds(0, 0, float64(world.width), float64(world.height)) {
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

func (o *Object) setId(id int) {
	o.Id = id
}

func (o *Object) getPos() physics.Point {
	return o.Pos
}

func (o *Object) getType() string {
	return o.Type
}
func (o *Object) getId() int {
	return o.Id
}
