package changelog

import (
	"aakimov/marsgame/physics"
	"aakimov/marsgame/world"
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
			ObjType: world.TypePlayer,
			ObjId:   11,
			Pos: &physics.Point{
				X: x,
				Y: y,
			},
			Angle:       &newR,
			CannonAngle: &newCR,
			Length:      &newL,
		}
		changeByTime.Add(changeByObject)

		changelog.AddAndCheckSize(changeByTime)
	}

	changelog.Optimize()
	require.Len(t, changelog.ChangesByTimeLog, 2)
	require.Equal(t, int64(0), changelog.ChangesByTimeLog[0].TimeId)
	require.Equal(t, int64(2), changelog.ChangesByTimeLog[1].TimeId)
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
			ObjType: world.TypePlayer,
			ObjId:   11,
			Pos: &physics.Point{
				X: x,
				Y: y,
			},
			Angle:       &newR,
			CannonAngle: &newCR,
			Length:      &newL,
		}
		changeByTime.Add(changeByObject)

		changelog.AddAndCheckSize(changeByTime)
	}

	changelog.Optimize()
	require.Len(t, changelog.ChangesByTimeLog, 3)
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
			ObjType: world.TypePlayer,
			ObjId:   11,
			Pos: &physics.Point{
				X: x,
				Y: y,
			},
			Angle:       &newR,
			CannonAngle: &newCR,
			Length:      &newL,
		}
		changeByTime.Add(changeByObject)

		changelog.AddAndCheckSize(changeByTime)
	}

	changelog.Optimize()
	require.Len(t, changelog.ChangesByTimeLog, 2)
}

func printMap(mapVar map[string]float64) {
	for k, v := range mapVar {
		fmt.Printf("%s = %f\n", k, v)
	}
}
