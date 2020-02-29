package physics

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDistance(t *testing.T) {
	p1 := &Point{X: 1, Y: 0}
	p2 := &Point{X: 10, Y: 0}
	require.Equal(t, 9., DistancePoints(p1, p2))
}

func TestAreObjsCollideFalse(t *testing.T) {
	o1 := &Obj{
		Pos:             Point{X: 1, Y: 0},
		CollisionRadius: 3,
	}
	o2 := &Obj{
		Pos:             Point{X: 9, Y: 0},
		CollisionRadius: 3,
	}
	require.False(t, AreObjsCollide(o1, o2))
}

func TestAreObjsCollideTrue(t *testing.T) {
	o1 := &Obj{
		Pos:             Point{X: 1, Y: 0},
		CollisionRadius: 5,
	}
	o2 := &Obj{
		Pos:             Point{X: 9, Y: 0},
		CollisionRadius: 5,
	}
	require.True(t, AreObjsCollide(o1, o2))
}

func BenchmarkDistance(b *testing.B) {
	p1 := &Point{X: 10.12, Y: 20.45}
	p2 := &Point{X: 53.12, Y: 21.45}
	for i := 0; i < b.N; i++ {
		_ = DistancePoints(p1, p2)
	}
}
