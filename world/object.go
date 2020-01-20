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
)

type IObject interface {
	run(world *World) *ChangeByObject
	setId(id int)
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
	changeByObject := ChangeByObject{
		ObjType: o.Type,
		ObjId:   strconv.Itoa(o.Id),
	}
	o.mu.Lock()
	if o.Speed != 0 {
		o.Pos.MoveForward(o.Angle, o.Speed)
		newPos := o.Pos
		newSpeed := o.Speed
		changeByObject.Pos = &newPos
		changeByObject.length = &newSpeed
		changeByObject.Angle = &o.Angle
	}
	o.mu.Unlock()

	if o.Speed != 0 {
		return &changeByObject
	}
	return nil
}

func (o *Object) setId(id int) {
	o.Id = id
}
