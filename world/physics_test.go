package world

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDistance(t *testing.T) {
	p1 := &Point{X: 1, Y: 0}
	p2 := &Point{X: 10, Y: 0}
	require.Equal(t, 9., distancePoints(p1, p2))
}

func TestAreObjsCollideFalse(t *testing.T) {
	o1 := &Object{
		Pos:             Point{X: 1, Y: 0},
		CollisionRadius: 3,
	}
	o2 := &Object{
		Pos:             Point{X: 9, Y: 0},
		CollisionRadius: 3,
	}
	require.False(t, areObjsCollide(o1, o2))
}

func TestAreObjsCollideTrue(t *testing.T) {
	o1 := &Object{
		Pos:             Point{X: 1, Y: 0},
		CollisionRadius: 5,
	}
	o2 := &Object{
		Pos:             Point{X: 9, Y: 0},
		CollisionRadius: 5,
	}
	require.True(t, areObjsCollide(o1, o2))
}

func BenchmarkDistance(b *testing.B) {
	p1 := &Point{X: 10.12, Y: 20.45}
	p2 := &Point{X: 53.12, Y: 21.45}
	for i := 0; i < b.N; i++ {
		_ = distancePoints(p1, p2)
	}
}
