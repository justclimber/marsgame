// Write Ahead Log - из мира баз данных
// wal содержит все изменения мира, предназначенные для обработки их клиентами -
// параллельными программами и браузерами
package wal

const defaultInt = 99999999
const defaultFloat = 99999999.
const LogBufferSize = 8

type Wal struct {
	curSize         int8
	logBuffer       *Log
	objectObservers []*ObjectObserver
	Sender          *Sender
}

func NewWal() *Wal {
	return &Wal{
		objectObservers: make([]*ObjectObserver, 0),
		logBuffer:       NewLog(),
		Sender:          NewSender(),
	}
}

type Log struct {
	TimeIds []int64
	Objects []*ObjectLog
}

func NewLog() *Log {
	return &Log{
		TimeIds: make([]int64, 0),
		Objects: make([]*ObjectLog, 0),
	}
}

func (w *Wal) Commit(timeId int64) {
	w.logBuffer.TimeIds = append(w.logBuffer.TimeIds, timeId)
	w.curSize++

	if w.curSize == LogBufferSize {
		w.Sender.logCh <- w.logBuffer
		w.curSize = 0
		w.resetLogBuffer()
	}
}

func (w *Wal) resetLogBuffer() {
	w.logBuffer = &Log{
		TimeIds: make([]int64, 0),
		Objects: make([]*ObjectLog, len(w.objectObservers)),
	}

	for i, oo := range w.objectObservers {
		ol := NewObjectLog(oo.Id, oo.ObjType)
		tl := NewTimeLog(false)
		oo.objectLog = ol
		oo.timeLog = tl
		oo.lastVelocityX = &tl.VelocityX
		oo.lastVelocityY = &tl.VelocityY
		oo.lastVelocityRotation = &tl.VelocityRotation
		oo.lastVelocityUntilTimeId = &tl.VelocityUntilTimeId
		oo.lastCannonRotation = &tl.CannonRotation
		oo.lastCannonUntilTimeId = &tl.CannonUntilTimeId

		w.logBuffer.Objects[i] = ol
	}
}

type ObjectLog struct {
	Id      uint32
	ObjType int8
	Times   []*TimeLog
}

func NewObjectLog(id uint32, objType int8) *ObjectLog {
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
	TimeId               int64
	IsNew                bool
	X                    float64
	Y                    float64
	Angle                float64
	CannonAngle          float64
	CannonRotation       float64
	CannonUntilTimeId    int64
	Fire                 bool
	Delete               bool
	Explode              bool
	ExplodeOther         bool
	DeleteOtherObjectIds []uint32
	VelocityX            float64
	VelocityY            float64
	VelocityRotation     float64
	VelocityUntilTimeId  int64
}

func NewTimeLog(isNew bool) *TimeLog {
	return &TimeLog{
		IsNew:               isNew,
		X:                   defaultInt,
		Y:                   defaultInt,
		Angle:               defaultFloat,
		CannonAngle:         defaultFloat,
		CannonRotation:      defaultFloat,
		CannonUntilTimeId:   defaultInt,
		VelocityX:           defaultFloat,
		VelocityY:           defaultFloat,
		VelocityRotation:    defaultFloat,
		VelocityUntilTimeId: defaultInt,
	}
}
