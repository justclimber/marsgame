package physics

func (p *Point) CheckIfOutOfBounds(x1, y1, x2, y2 float64) bool {
	return p.X < x1 || p.Y < y1 || p.X > x2 || p.Y > y2
}
