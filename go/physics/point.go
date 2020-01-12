package physics

import "math"

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (p *Point) MoveForward(angle float64, length float64) {
	p.X += math.Cos(angle) * length
	p.Y += math.Sin(angle) * length
}
