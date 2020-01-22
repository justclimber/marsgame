package world

import "math"

type Point struct {
	x float64
	y float64
}

func (p *Point) MoveForward(angle float64, length float64) *Vector {
	dx := math.Cos(angle) * length
	dy := math.Sin(angle) * length
	p.x += dx
	p.y += dy

	return &Vector{Point{x: dx, y: dy}}
}

type Vector struct {
	Point
}
