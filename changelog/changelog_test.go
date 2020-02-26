package changelog

import (
	"aakimov/marsgame/physics"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChangeLog_Marshal(t *testing.T) {
	x := 0.
	y := 2.
	r := 5.12312312
	cr := 10.
	l := 1.5

	changelog := NewChangeLog()
	changeByTime := NewChangeByTime(int64(1))
	changeByObject := &ChangeByObject{
		ObjType: "player",
		ObjId:   11,
		Pos: &physics.Point{
			X: x,
			Y: y,
		},
		Angle:       &r,
		CannonAngle: &cr,
		Length:      &l,
	}
	changeByTime.Add(changeByObject)

	changelog.AddAndCheckSize(changeByTime)

	jsonOut, err := json.Marshal(changelog.ChangesByTimeLog)
	assert.Nil(t, err)
	require.Equal(t, `[{"tId":1,"chObjs":[{"t":"player","id":11,"x":0,"y":2,"a":5.12,"ca":10}]}]`, string(jsonOut))
}
