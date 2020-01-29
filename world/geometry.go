package world

import "math"

type Point struct {
	X float64
	Y float64
}

func (p *Point) MoveForward(angle float64, length float64) *Vector {
	dx := math.Cos(angle) * length
	dy := math.Sin(angle) * length
	p.X += dx
	p.Y += dy

	return &Vector{dx, dy}
}

func (p *Point) add(v1 *Vector) *Point {
	return &Point{p.X + v1.X, p.Y + v1.Y}
}

type Vector struct {
	X float64
	Y float64
}

func (v *Vector) add(v1 *Vector) *Vector {
	return &Vector{v.X + v1.X, v.Y + v1.Y}
}

func (v *Vector) multiplyOnScalar(k float64) *Vector {
	return &Vector{v.X * k, v.Y * k}
}

func (v *Vector) multiplyOnVector(v1 *Vector) float64 {
	return v.X*v1.X + v.Y*v1.Y
}

func (v *Vector) len() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func makeNormalVectorByAngle(angle float64) *Vector {
	return &Vector{
		X: math.Cos(angle),
		Y: math.Sin(angle),
	}
}
