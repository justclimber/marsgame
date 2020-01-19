package world

import "aakimov/marsgame/physics"

type Object struct {
	Pos    physics.Point
	Width  int
	Height int
	Angle  float64
}
