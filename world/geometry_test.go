package world

import (
	"github.com/stretchr/testify/require"

	"math"
	"testing"
)

func TestAngle(t *testing.T) {
	require.Equal(t, 0., angle(5., 5., 15., 5.), "0")
	require.Equal(t, math.Pi/4, angle(-5., 2., 0., 7.), "45")
	require.Equal(t, math.Pi/2, angle(-5., -5., -5., -1.), "90")
	require.Equal(t, math.Pi, angle(5., -5., -5., -5.), "180")
	require.Equal(t, math.Pi+math.Pi/4, angle(5., -5., 0., -10.), "180-45")
}

func TestAngleOrigin(t *testing.T) {
	require.Equal(t, 0., angleOrigin(0., 0.), "0")
	require.Equal(t, math.Pi/4, angleOrigin(10., 10.), "45")
	require.Equal(t, math.Pi/2, angleOrigin(0., 10.), "90")
	require.Equal(t, math.Pi, angleOrigin(-10., 0.), "180")
	require.Equal(t, math.Pi+math.Pi/4, angleOrigin(-10., -10.), "180+45")
	require.Equal(t, math.Pi+math.Pi/2+math.Pi/4, angleOrigin(10., -10.), "180+90+45")

	a := angleOrigin(10., -1.)
	near := a < math.Pi*2 && a > math.Pi*2-0.1
	require.True(t, near, "~360")
}

func TestAngleOriginByCircle(t *testing.T) {
	var x, y float64
	r := 10.
	lastAngle := 0.
	for a := 0.; a < math.Pi*2; a += 0.1 {
		x = r * math.Cos(a)
		y = r * math.Sin(a)
		resultAngle := angleOrigin(x, y)
		require.LessOrEqualf(t, lastAngle, resultAngle, "a: %f, r: %f", a, resultAngle)
		lastAngle = resultAngle
	}
}
