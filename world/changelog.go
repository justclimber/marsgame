package world

import (
	"encoding/json"
	"math"
)

const ChangelogChannelBufferSize = 10
const ChangelogBufferSize = 8

type ChangeLog struct {
	changesByTimeCh  chan *ChangeByTime
	changesByTimeLog []*ChangeByTime
}

type ChangeByTime struct {
	TimeId          int64             `json:"tId"`
	ChangesByObject []*ChangeByObject `json:"chObjs"`
}

type ChangeByObject struct {
	ObjType     string
	ObjId       string
	Pos         *Point
	Angle       *float64
	CannonAngle *float64
	Delete      bool
	length      *float64
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
		changesByTimeCh:  make(chan *ChangeByTime, ChangelogChannelBufferSize),
		changesByTimeLog: make([]*ChangeByTime, 0, ChangelogBufferSize),
	}
}

func (ch *ChangeLog) AddToBuffer(changeByTime *ChangeByTime) {
	ch.changesByTimeCh <- changeByTime
}

func (ch *ChangeLog) AddAndCheckSize(changeByTime *ChangeByTime) bool {
	ch.changesByTimeLog = append(ch.changesByTimeLog, changeByTime)
	return len(ch.changesByTimeLog) >= ChangelogBufferSize
}

func (ch *ChangeByObject) MarshalJSON() ([]byte, error) {
	var xp, yp *int
	if ch.Pos != nil {
		x, y := int(ch.Pos.x), int(ch.Pos.y)
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
		ObjType     string   `json:"t"`
		ObjId       string   `json:"id"`
		X           *int     `json:"x,omitempty"`
		Y           *int     `json:"y,omitempty"`
		Angle       *float64 `json:"a,omitempty"`
		CannonAngle *float64 `json:"ca,omitempty"`
		Delete      *bool    `json:"d,omitempty"`
	}{
		ObjType:     ch.ObjType,
		ObjId:       ch.ObjId,
		X:           xp,
		Y:           yp,
		Angle:       angle,
		CannonAngle: cannonAngle,
		Delete:      dp,
	})
}
