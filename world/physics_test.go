package world

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDistance(t *testing.T) {
	p1 := &Point{x: 1, y: 0}
	p2 := &Point{x: 10, y: 0}
	require.Equal(t, 9., distance(p1, p2))
}

func TestAreObjsCollideFalse(t *testing.T) {
	o1 := &Object{
		Pos:             Point{x: 1, y: 0},
		CollisionRadius: 3,
	}
	o2 := &Object{
		Pos:             Point{x: 9, y: 0},
		CollisionRadius: 3,
	}
	require.False(t, areObjsCollide(o1, o2))
}

func TestAreObjsCollideTrue(t *testing.T) {
	o1 := &Object{
		Pos:             Point{x: 1, y: 0},
		CollisionRadius: 5,
	}
	o2 := &Object{
		Pos:             Point{x: 9, y: 0},
		CollisionRadius: 5,
	}
	require.True(t, areObjsCollide(o1, o2))
}

func BenchmarkDistance(b *testing.B) {
	p1 := &Point{x: 10.12, y: 20.45}
	p2 := &Point{x: 53.12, y: 21.45}
	for i := 0; i < b.N; i++ {
		_ = distance(p1, p2)
	}
}
