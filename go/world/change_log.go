package world

import (
	"aakimov/marsgame/go/physics"
	"aakimov/marsgame/go/server"
	"encoding/json"
	"fmt"
	"log"
	"math"
)

const (
	TypePlayer = "player"
	TypeObject = "object"
)

const ChangelogChannelBufferSize = 10
const ChangelogBufferSize = 8

type ChangeLog struct {
	changesByTimeCh  chan *ChangeByTime
	changesByTimeLog []*ChangeByTime
}

type ChangeByTime struct {
	TimeId          int64             `json:"timeId"`
	ChangesByObject []*ChangeByObject `json:"changesByObject"`
}

type ChangeByObject struct {
	ObjType string        `json:"objType"`
	ObjId   string        `json:"objId"`
	Pos     physics.Point `json:"pos"`
	Angle   float64       `json:"angle"`
	length  float64
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

// Optimize intermediate changelog for individual object if they can be interpolated (have const diffs fo 3 changes)
func (ch *ChangeLog) Optimize() {
	changelogLen := len(ch.changesByTimeLog)
	if changelogLen < 3 {
		return
	}
	i1, i2, ok := lookupForSameDiff(changelogLen-1, changelogLen-2, ch)
	if ok {
		log.Println("Optimized!")
		ch.cutInterpolableChanges(i1, i2)
	}
}

func (ch *ChangeLog) cutInterpolableChanges(i1, i2 int) {
	fmt.Printf("Cut from %d, %d, len: %d\n", i2, i1, len(ch.changesByTimeLog))
	ch.changesByTimeLog = append(ch.changesByTimeLog[:i2+1], ch.changesByTimeLog[i1:]...)
}

func lookupForSameDiff(tailIndex int, index int, ch *ChangeLog) (int, int, bool) {
	dl, dr := getDiff(tailIndex, index, ch)
	dl1, dr1 := getDiff(tailIndex-1, index-1, ch)

	if checkAreDiffEqualZero(dl, dl1) && checkAreDiffTheSame(dr, dr1) {
		if index > 2 {
			_, index1, ok := lookupForSameDiff(tailIndex, index-2, ch)
			if ok {
				return tailIndex, index1, ok
			}
			//fmt.Println(tailIndex, index - 1)
			//fmt.Println(len(ch.changesByTimeLog))
			//fmt.Println(len(ch.changesByTimeLog[6:]))
			//ch.cutInterpolableChanges(tailIndex, index - 1)
			return 0, 0, false
			//if index > 5 {
			//	lookupForSameDiff(tailIndex - 2, index - 3, ch)
			//}
		}
		//ch.cutInterpolableChanges(tailIndex, index - 1)
		return tailIndex, index - 1, true
		//return 0, 0, false
	}
	return 0, 0, false
}

const floatDelta = 0.01

func checkAreDiffTheSame(d1, d2 map[string]float64) bool {
	for k, v := range d1 {
		if d2[k]-v > floatDelta {
			return false
		}
	}
	return true
}

func checkAreDiffEqualZero(d1, d2 map[string]float64) bool {
	for k, v := range d1 {
		if math.Abs(d2[k]) > floatDelta || math.Abs(v) > floatDelta {
			return false
		}
	}
	return true
}

func getDiff(index int, index1 int, ch *ChangeLog) (map[string]float64, map[string]float64) {
	l, r := getValuesForChanges(ch.changesByTimeLog[index])
	l1, r1 := getValuesForChanges(ch.changesByTimeLog[index1])
	dl := make(map[string]float64)
	dr := make(map[string]float64)

	for k, v := range l {
		dl[k] = v - l1[k]
	}
	for k, v := range r {
		dr[k] = v - r1[k]
	}
	return dl, dr
}

func getValuesForChanges(changeByTime *ChangeByTime) (map[string]float64, map[string]float64) {
	l := make(map[string]float64)
	r := make(map[string]float64)
	for _, changeByObject := range changeByTime.ChangesByObject {
		key := changeByObject.ObjId + changeByObject.ObjType
		l[key] = changeByObject.length
		r[key] = changeByObject.Angle
	}
	return l, r
}

func PackChangesToCommand(changes []*ChangeByTime) *server.Command {
	command := server.Command{Type: "worldChanges"}
	jsonBytes, err := json.Marshal(changes)
	if err != nil {
		log.Println(err)
	}
	command.Payload = string(jsonBytes)
	return &command
}
