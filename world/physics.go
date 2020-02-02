package world

func (o *Object) isCollideWith(o1 IObject) bool {
	p := o1.getPos()
	return distancePoints(&o.Pos, &p) <= float64(o.CollisionRadius+o1.getCollisionRadius())
}

func areObjsCollide(o1, o2 *Object) bool {
	return distancePoints(&o1.Pos, &o2.Pos) <= float64(o1.CollisionRadius+o2.CollisionRadius)
}
