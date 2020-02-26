package changelog

import (
	"aakimov/marsgame/physics"
	"encoding/json"
	"math"
)

const ChannelBufferSize = 10
const BufferSize = 8

type ChangeLog struct {
	ChangesByTimeCh  chan *ChangeByTime
	TerminateCh      chan bool
	ChangesByTimeLog []*ChangeByTime
}

type ChangeByTime struct {
	TimeId          int64             `json:"tId"`
	ChangesByObject []*ChangeByObject `json:"chObjs"`
}

type ChangeByObject struct {
	ObjType       string
	ObjId         uint32
	Pos           *physics.Point
	Angle         *float64
	CannonAngle   *float64
	Delete        bool
	DeleteOtherId int
	Length        *float64
}

func NewChangeByTime(timeId int64) *ChangeByTime {
	return &ChangeByTime{
		TimeId:          timeId,
		ChangesByObject: make([]*ChangeByObject, 0),
	}
}

func (ch *ChangeByTime) Add(changeByObject *ChangeByObject) {
	ch.ChangesByObject = append(ch.ChangesByObject, changeByObject)
}

func (ch *ChangeByTime) IsNotEmpty() bool {
	return len(ch.ChangesByObject) > 0
}

func NewChangeLog() *ChangeLog {
	return &ChangeLog{
		ChangesByTimeCh:  make(chan *ChangeByTime, ChannelBufferSize),
		TerminateCh:      make(chan bool, 1),
		ChangesByTimeLog: make([]*ChangeByTime, 0, BufferSize),
	}
}

func (ch *ChangeLog) Terminate() {
	ch.TerminateCh <- true
}

func (ch *ChangeLog) Reset() {
	ch.ChangesByTimeLog = make([]*ChangeByTime, 0, BufferSize)
}

func (ch *ChangeLog) GetLog() []*ChangeByTime {
	return ch.ChangesByTimeLog
}

func (ch *ChangeLog) AddToBuffer(changeByTime *ChangeByTime) {
	ch.ChangesByTimeCh <- changeByTime
}

func (ch *ChangeLog) AddAndCheckSize(changeByTime *ChangeByTime) bool {
	ch.ChangesByTimeLog = append(ch.ChangesByTimeLog, changeByTime)
	return len(ch.ChangesByTimeLog) >= BufferSize
}

func (ch *ChangeByObject) MarshalJSON() ([]byte, error) {
	var xp, yp *int
	if ch.Pos != nil {
		x, y := int(ch.Pos.X), int(ch.Pos.Y)
		xp, yp = &x, &y
	}
	var angle, cannonAngle *float64
	if ch.Angle != nil {
		a := math.Round(*ch.Angle*100) / 100
		angle = &a
	}
	if ch.CannonAngle != nil {
		ca := math.Round(*ch.CannonAngle*100) / 100
		cannonAngle = &ca
	}
	var dp *bool
	if ch.Delete {
		dp = &ch.Delete
	}
	return json.Marshal(struct {
		ObjType       string   `json:"t"`
		ObjId         uint32   `json:"id"`
		X             *int     `json:"x,omitempty"`
		Y             *int     `json:"y,omitempty"`
		Angle         *float64 `json:"a,omitempty"`
		CannonAngle   *float64 `json:"ca,omitempty"`
		Delete        *bool    `json:"d,omitempty"`
		DeleteOtherId int      `json:"did,omitempty"`
	}{
		ObjType:       ch.ObjType,
		ObjId:         ch.ObjId,
		X:             xp,
		Y:             yp,
		Angle:         angle,
		CannonAngle:   cannonAngle,
		Delete:        dp,
		DeleteOtherId: ch.DeleteOtherId,
	})
}
