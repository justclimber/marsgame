package world

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChangeLog_OptimizeSimplePositiveCase(t *testing.T) {
	dx := 3.
	dy := 2.
	dr := 1.
	dcr := 1.
	x := 0.
	y := 0.
	r := 0.
	cr := 0.
	l := 1.5

	changelog := NewChangeLog()
	for i := 0; i <= 2; i++ {
		changeByTime := NewChangeByTime(int64(i))

		x = x + dx
		y = y + dy
		r = r + dr
		cr = cr + dcr
		newL := l
		newR := r
		newCR := cr

		changeByObject := &ChangeByObject{
			ObjType: TypePlayer,
			ObjId:   "11",
			Pos: &Point{
				x: x,
				y: y,
			},
			Angle:       &newR,
			CannonAngle: &newCR,
			length:      &newL,
		}
		changeByTime.Add(changeByObject)

		changelog.AddAndCheckSize(changeByTime)
	}

	changelog.Optimize()
	require.Len(t, changelog.changesByTimeLog, 2)
	require.Equal(t, int64(0), changelog.changesByTimeLog[0].TimeId)
	require.Equal(t, int64(2), changelog.changesByTimeLog[1].TimeId)
}

func TestChangeLog_OptimizeSimpleNegative(t *testing.T) {
	dx := 3.
	dy := 2.
	dr := 1.
	dcr := 1.
	x := 0.
	y := 0.
	r := 0.
	cr := 0.
	l := 1.5

	changelog := NewChangeLog()
	for i := 0; i <= 2; i++ {
		changeByTime := NewChangeByTime(int64(i))

		x = x + dx
		y = y + dy
		r = r + dr
		cr = cr + dcr

		// make some diff shake here
		if i == 1 {
			l = l + 5.
		}
		newL := l
		newR := r
		newCR := cr

		changeByObject := &ChangeByObject{
			ObjType: TypePlayer,
			ObjId:   "11",
			Pos: &Point{
				x: x,
				y: y,
			},
			Angle:       &newR,
			CannonAngle: &newCR,
			length:      &newL,
		}
		changeByTime.Add(changeByObject)

		changelog.AddAndCheckSize(changeByTime)
	}

	changelog.Optimize()
	require.Len(t, changelog.changesByTimeLog, 3)
}

func TestChangeLog_OptimizePositiveMoreElements(t *testing.T) {
	dx := 3.
	dy := 2.
	dr := 1.
	dcr := 1.
	x := 0.
	y := 0.
	r := 0.
	cr := 0.
	l := 1.5

	changelog := NewChangeLog()
	for i := 0; i <= 6; i++ {
		changeByTime := NewChangeByTime(int64(i))

		x = x + dx
		y = y + dy
		r = r + dr
		cr = cr + dcr
		newL := l
		newR := r
		newCR := cr

		changeByObject := &ChangeByObject{
			ObjType: TypePlayer,
			ObjId:   "11",
			Pos: &Point{
				x: x,
				y: y,
			},
			Angle:       &newR,
			CannonAngle: &newCR,
			length:      &newL,
		}
		changeByTime.Add(changeByObject)

		changelog.AddAndCheckSize(changeByTime)
	}

	changelog.Optimize()
	require.Len(t, changelog.changesByTimeLog, 2)
}

func printMap(mapVar map[string]float64) {
	for k, v := range mapVar {
		fmt.Printf("%s = %f\n", k, v)
	}
}
