package world

import "math"

type Point struct {
	X float64
	Y float64
}

func (p *Point) checkIfOutOfBounds(x1, y1, x2, y2 float64) bool {
	return p.X < x1 || p.Y < y1 || p.X > x2 || p.Y > y2
}

func (p *Point) distanceTo(p2 *Point) float64 {
	dx := p.X - p2.X
	dy := p.Y - p2.Y
	ds := (dx * dx) + (dy * dy)

	return math.Sqrt(ds)
}

func distancePoints(p1, p2 *Point) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	ds := dx*dx + dy*dy

	return math.Sqrt(ds)
}

func distance(x1, y1, x2, y2 float64) float64 {
	dx := x1 - x2
	dy := y1 - y2
	ds := (dx * dx) + (dy * dy)

	return math.Sqrt(ds)
}

func angle(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1

	return angleOrigin(dx, dy)
}

func angleOrigin(x, y float64) float64 {
	atan := math.Atan(y / x)
	if math.IsNaN(atan) {
		return 0.
	}
	switch {
	case x < 0 && y >= 0:
		return math.Pi + atan
	case x < 0 && y < 0:
		return math.Pi + atan
	case x > 0 && y < 0:
		return math.Pi*2 + atan
	default:
		return atan
	}
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
