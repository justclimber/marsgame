package world

import (
	"aakimov/marsgame/backend/physics"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChangeLog_Marshal(t *testing.T) {
	x := 0.
	y := 0.
	r := 0.
	cr := 0.
	l := 1.5

	changelog := NewChangeLog()
	changeByTime := NewChangeByTime(int64(1))
	changeByObject := &ChangeByObject{
		ObjType: TypePlayer,
		ObjId:   "11",
		Pos: &physics.Point{
			X: x,
			Y: y,
		},
		Angle:       &r,
		CannonAngle: &cr,
		length:      &l,
	}
	changeByTime.Add(changeByObject)

	changelog.AddAndCheckSize(changeByTime)

	_, err := json.Marshal(changelog.changesByTimeLog)
	assert.Nil(t, err)
	//require.Equal(t, "asd", string(jsonOut))
}
