package world

import (
	"math"
)

func (p *Point) checkIfOutOfBounds(x1, y1, x2, y2 float64) bool {
	return p.X < x1 || p.Y < y1 || p.X > x2 || p.Y > y2
}

func (p *Point) distanceTo(p2 *Point) float64 {
	dx := p.X - p2.X
	dy := p.Y - p2.Y
	ds := (dx * dx) + (dy * dy)

	return math.Sqrt(ds)
}

func distance(p1, p2 *Point) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	ds := dx*dx + dy*dy

	return math.Sqrt(ds)
}

func (o *Object) isCollideWith(o1 IObject) bool {
	p := o1.getPos()
	return distance(&o.Pos, &p) <= float64(o.CollisionRadius+o1.getCollisionRadius())
}

func areObjsCollide(o1, o2 *Object) bool {
	return distance(&o1.Pos, &o2.Pos) <= float64(o1.CollisionRadius+o2.CollisionRadius)
}
