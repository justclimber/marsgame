// Write Ahead Log - из мира баз данных
// wal содержит все изменения мира, предназначенные для обработки их клиентами -
// параллельными программами и браузерами
package wal

type Wal struct {
	MainLog        *Log
	objectManagers []*ObjectObserver
	Sender         *Sender
}

func NewWal() *Wal {
	return &Wal{
		objectManagers: make([]*ObjectObserver, 0),
		MainLog:        NewLog(),
		Sender:         NewSender(),
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

func (w *Wal) Commit(timeId int64) {
	w.MainLog.TimeIds = append(w.MainLog.TimeIds, timeId)
	//w.Sender.logCh <- w.MainLog
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
