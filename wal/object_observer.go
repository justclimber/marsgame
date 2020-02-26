package wal

import "aakimov/marsgame/physics"

type ObjectObserver struct {
	Id                      int
	ObjType                 string
	timeLog                 *TimeLog
	objectLog               *ObjectLog
	lastVelocityX           *float64
	lastVelocityY           *float64
	lastVelocityRotation    *float64
	lastVelocityUntilTimeId *int64
	lastCannonRotation      *float64
	lastCannonUntilTimeId   *int64
}

func (w *Wal) NewObjectObserver(id int, objType string) *ObjectObserver {
	ol := NewObjectLog(id, objType)
	tl := NewTimeLog()
	o := &ObjectObserver{
		Id:        id,
		ObjType:   objType,
		objectLog: ol,
		timeLog:   tl,
	}

	o.lastVelocityX = tl.VelocityX
	o.lastVelocityY = tl.VelocityY
	o.lastVelocityRotation = tl.VelocityRotation
	o.lastVelocityUntilTimeId = tl.VelocityUntilTimeId
	o.lastCannonRotation = tl.CannonRotation
	o.lastCannonUntilTimeId = tl.CannonUntilTimeId

	w.objectManagers = append(w.objectManagers, o)
	w.MainLog.Objects = append(w.MainLog.Objects, ol)
	return o
}

func (oo *ObjectObserver) AddRotation(rotation float64) {
	oo.timeLog.VelocityRotation = &rotation
}

func (oo *ObjectObserver) AddAngle(angle float64) {
	oo.timeLog.Angle = &angle
}

func (oo *ObjectObserver) AddPosAndVelocity(pos physics.Point, velocity *physics.Vector) {
	oo.timeLog.X = &pos.X
	oo.timeLog.Y = &pos.Y
	oo.timeLog.VelocityX = &velocity.X
	oo.timeLog.VelocityY = &velocity.Y
}

func (oo *ObjectObserver) AddCannonRotation(rotation float64) {
	oo.timeLog.CannonRotation = &rotation
}

func (oo *ObjectObserver) AddCannonAngle(angle float64) {
	oo.timeLog.CannonAngle = &angle
}

func (oo *ObjectObserver) AddShoot() {
	fire := true
	oo.timeLog.Fire = &fire
}

func (oo *ObjectObserver) AddDelete() {
	del := true
	oo.timeLog.Delete = &del
}

func (oo *ObjectObserver) AddExplode() {
	explode := true
	oo.timeLog.Explode = &explode
}

func (oo *ObjectObserver) AddDeleteOtherIds(ids []int) {
	oo.timeLog.DeleteOtherObjectIds = ids
}

func (oo *ObjectObserver) Commit(timeId int64) {
	oo.timeLog.TimeId = timeId
	oo.optimize()
	if !oo.timeLog.skip {
		oo.objectLog.Times = append(oo.objectLog.Times, oo.timeLog)
	}
	oo.timeLog = NewTimeLog()
}
