package wal

import (
	"github.com/justclimber/marsgame/helpers"
	"github.com/justclimber/marsgame/physics"
)

type ObjectObserver struct {
	Id                      uint32
	ObjType                 int8
	timeLog                 *TimeLog
	objectLog               *ObjectLog
	lastVelocityLen         *float64
	lastVelocityRotation    *float64
	lastVelocityUntilTimeId *int64
	lastCannonRotation      *float64
	lastCannonUntilTimeId   *int64
	toDelete                bool
}

func (w *Wal) NewObjectObserver(id uint32, objType int8) *ObjectObserver {
	ol := NewObjectLog(id, objType)
	tl := NewTimeLog(true)
	o := &ObjectObserver{
		Id:        id,
		ObjType:   objType,
		objectLog: ol,
		timeLog:   tl,
	}

	o.lastVelocityLen = &tl.VelocityLen
	o.lastVelocityRotation = &tl.VelocityRotation
	o.lastVelocityUntilTimeId = &tl.VelocityUntilTimeId
	o.lastCannonRotation = &tl.CannonRotation
	o.lastCannonUntilTimeId = &tl.CannonUntilTimeId

	w.objectObservers[id] = o
	w.logBuffer.Objects[id] = ol
	return o
}

func (oo *ObjectObserver) AddRotation(rotation float64) {
	oo.timeLog.VelocityRotation = helpers.Round(rotation)
}

func (oo *ObjectObserver) AddAngle(angle float64) {
	oo.timeLog.Angle = angle
}

func (oo *ObjectObserver) AddPosAndVelocityLen(pos physics.Point, velocityLen float64) {
	oo.timeLog.X = pos.X
	oo.timeLog.Y = pos.Y
	oo.timeLog.VelocityLen = helpers.Round(velocityLen)
}

func (oo *ObjectObserver) AddCannonRotation(rotation float64) {
	oo.timeLog.CannonRotation = rotation
}

func (oo *ObjectObserver) AddCannonAngle(angle float64) {
	oo.timeLog.CannonAngle = angle
}

func (oo *ObjectObserver) AddShoot() {
	oo.timeLog.Fire = true
}

func (oo *ObjectObserver) AddDelete() {
	oo.timeLog.Delete = true
	oo.toDelete = true
}

func (oo *ObjectObserver) AddExplode() {
	oo.timeLog.Delete = true
	oo.timeLog.Explode = true
}

func (oo *ObjectObserver) AddDeleteOtherIds(ids []uint32) {
	oo.timeLog.DeleteOtherObjectIds = ids
}

func (oo *ObjectObserver) Commit(timeId int64) {
	oo.timeLog.TimeId = timeId
	oo.optimize()
	if !oo.timeLog.skip {
		oo.objectLog.Times = append(oo.objectLog.Times, oo.timeLog)
	}
	oo.timeLog = NewTimeLog(false)
}
