// Write Ahead Log - из мира баз данных
// wal содержит все изменения мира, предназначенные для обработки их клиентами -
// параллельными программами и браузерами
package wal

const defaultInt = 99999999
const defaultFloat = 99999999
const LogBufferSize = 8

type Wal struct {
	curSize         int8
	logBuffer       *Log
	objectObservers map[uint32]*ObjectObserver
	Sender          *Sender
}

func NewWal() *Wal {
	storage := NewLog()
	return &Wal{
		objectObservers: make(map[uint32]*ObjectObserver),
		logBuffer:       NewLog(),
		Sender:          NewSender(storage),
	}
}

type Log struct {
	currTimeId int64
	TimeIds    []int64
	Objects    map[uint32]*ObjectLog
}

func NewLog() *Log {
	return &Log{
		TimeIds: make([]int64, 0),
		Objects: make(map[uint32]*ObjectLog),
	}
}

func (l *Log) merge(withLog *Log) {
	l.currTimeId = withLog.currTimeId
	l.TimeIds = append(l.TimeIds, withLog.TimeIds...)
	for k, v := range withLog.Objects {
		if _, ok := l.Objects[k]; ok {
			l.Objects[k].merge(v)
		} else {
			l.Objects[k] = v
		}
	}
}

func (w *Wal) Commit(timeId int64) {
	w.logBuffer.TimeIds = append(w.logBuffer.TimeIds, timeId)
	w.logBuffer.currTimeId = timeId
	w.curSize++

	if w.curSize == LogBufferSize {
		for k, o := range w.logBuffer.Objects {
			if len(o.Times) == 0 {
				delete(w.logBuffer.Objects, k)
			}
		}
		//helpers.PrettyPrint("wal", w.logBuffer)
		w.Sender.logCh <- w.logBuffer
		for k, oo := range w.objectObservers {
			if oo.toDelete {
				delete(w.objectObservers, k)
				delete(w.logBuffer.Objects, k)
			}
		}
		w.curSize = 0
		w.resetLogBuffer()
	}
}

func (w *Wal) resetLogBuffer() {
	w.logBuffer = &Log{
		TimeIds: make([]int64, 0),
		Objects: make(map[uint32]*ObjectLog),
	}

	w.logBuffer.Objects = make(map[uint32]*ObjectLog, len(w.objectObservers))
	for _, oo := range w.objectObservers {
		ol := NewObjectLog(oo.Id, oo.ObjType)
		tl := NewTimeLog(false)
		oo.objectLog = ol
		oo.timeLog = tl
		oo.lastVelocityLen = &tl.VelocityLen
		oo.lastVelocityRotation = &tl.VelocityRotation
		oo.lastVelocityUntilTimeId = &tl.VelocityUntilTimeId
		oo.lastCannonRotation = &tl.CannonRotation
		oo.lastCannonUntilTimeId = &tl.CannonUntilTimeId
		w.logBuffer.Objects[oo.Id] = ol
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

func (o *ObjectLog) merge(withObjectLog *ObjectLog) {
	o.Times = append(o.Times, withObjectLog.Times...)
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
	VelocityLen          float64
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
		VelocityLen:         defaultFloat,
		VelocityRotation:    defaultFloat,
		VelocityUntilTimeId: defaultInt,
	}
}
