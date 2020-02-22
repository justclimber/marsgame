// Write Ahead Log - из мира баз данных
// wal содержит все изменения мира, предназначенные для обработки их клиентами -
// параллельными программами и браузерами
package wal

import (
	"aakimov/marsgame/physics"
)

type Manager struct {
	MainLog        *Log
	objectManagers []*ObjectManager
	Sender         *WallSender
}

func NewManager() *Manager {
	return &Manager{
		objectManagers: make([]*ObjectManager, 0),
		MainLog:        NewLog(),
		Sender:         NewWallSender(),
	}
}

type Log struct {
	TimeIds []int64      `json:"tIds"`
	Objects []*ObjectLog `json:"objs"`
}

func NewLog() *Log {
	return &Log{
		TimeIds: make([]int64, 0),
		Objects: make([]*ObjectLog, 0),
	}
}

func (l *Manager) CreateObjectManager(id int, objType string) *ObjectManager {
	ol := NewObjectLog(id, objType)
	tl := NewTimeLog()
	o := &ObjectManager{
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

	l.objectManagers = append(l.objectManagers, o)
	l.MainLog.Objects = append(l.MainLog.Objects, ol)
	return o
}

func (l *Manager) Commit(timeId int64) {
	l.MainLog.TimeIds = append(l.MainLog.TimeIds, timeId)
	//l.Sender.logCh <- l.MainLog
}

type ObjectManager struct {
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

type ObjectLog struct {
	Id      int        `json:"i"`
	ObjType string     `json:"t"`
	Times   []*TimeLog `json:"l"`
}

func NewObjectLog(id int, objType string) *ObjectLog {
	return &ObjectLog{
		Id:      id,
		ObjType: objType,
		Times:   make([]*TimeLog, 0),
	}
}

func (o *ObjectLog) LastTimeLog() *TimeLog {
	if len(o.Times) == 0 {
		return nil
	}
	return o.Times[len(o.Times)-1]
}

type TimeLog struct {
	skip                 bool
	TimeId               int64    `json:"i"`
	IsNew                *bool    `json:"n,omitempty"`
	X                    *float64 `json:"x,omitempty"`
	Y                    *float64 `json:"y,omitempty"`
	Angle                *float64 `json:"a,omitempty"`
	CannonAngle          *float64 `json:"ca,omitempty"`
	CannonRotation       *float64 `json:"cr,omitempty"`
	CannonUntilTimeId    *int64   `json:"cu,omitempty"`
	Fire                 *bool    `json:"f,omitempty"`
	Delete               *bool    `json:"d,omitempty"`
	Explode              *bool    `json:"e,omitempty"`
	ExplodeOther         *bool    `json:"eo,omitempty"`
	DeleteOtherObjectIds []int    `json:"did,omitempty"`
	VelocityX            *float64 `json:"vx,omitempty"`
	VelocityY            *float64 `json:"vy,omitempty"`
	VelocityRotation     *float64 `json:"vr,omitempty"`
	VelocityUntilTimeId  *int64   `json:"vu,omitempty"`
}

func NewTimeLog() *TimeLog {
	return &TimeLog{}
}

func (om *ObjectManager) AddRotation(rotation float64) {
	om.timeLog.VelocityRotation = &rotation
}

func (om *ObjectManager) AddAngle(angle float64) {
	om.timeLog.Angle = &angle
}

func (om *ObjectManager) AddPosAndVelocity(pos physics.Point, velocity *physics.Vector) {
	om.timeLog.X = &pos.X
	om.timeLog.Y = &pos.Y
	om.timeLog.VelocityX = &velocity.X
	om.timeLog.VelocityY = &velocity.Y
}

func (om *ObjectManager) AddCannonRotation(rotation float64) {
	om.timeLog.CannonRotation = &rotation
}

func (om *ObjectManager) AddCannonAngle(angle float64) {
	om.timeLog.CannonAngle = &angle
}

func (om *ObjectManager) AddShoot() {
	fire := true
	om.timeLog.Fire = &fire
}

func (om *ObjectManager) AddDelete() {
	del := true
	om.timeLog.Delete = &del
}

func (om *ObjectManager) AddExplode() {
	explode := true
	om.timeLog.Explode = &explode
}

func (om *ObjectManager) AddDeleteOtherIds(ids []int) {
	om.timeLog.DeleteOtherObjectIds = ids
}

func (om *ObjectManager) Commit(timeId int64) {
	om.timeLog.TimeId = timeId
	om.optimize()
	if !om.timeLog.skip {
		om.objectLog.Times = append(om.objectLog.Times, om.timeLog)
	}
	om.timeLog = NewTimeLog()
}
