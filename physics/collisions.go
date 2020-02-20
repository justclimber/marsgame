package physics

func (o *Obj) IsCollideWith(o1 Obj) bool {
	p := o1.Pos
	return DistancePoints(&o.Pos, &p) <= float64(o.CollisionRadius+o1.CollisionRadius)
}

func AreObjsCollide(o1, o2 *Obj) bool {
	return DistancePoints(&o1.Pos, &o2.Pos) <= float64(o1.CollisionRadius+o2.CollisionRadius)
}
