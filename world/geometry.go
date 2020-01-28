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

	return &Vector{dx, dy}
}

func (p *Point) add(v1 Vector) Point {
	return Point{p.x + v1.x, p.y + v1.y}
}

type Vector struct {
	x float64
	y float64
}

func (v *Vector) add(v1 Vector) Vector {
	return Vector{v.x + v1.x, v.y + v1.y}
}

func (v *Vector) multiplyOnScalar(k float64) Vector {
	return Vector{v.x * k, v.y * k}
}

func (v *Vector) multiplyOnVector(v1 Vector) float64 {
	return v.x*v1.x + v.y*v1.y
}

func (v *Vector) len() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}
