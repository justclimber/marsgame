package world

import (
	"aakimov/marsgame/physics"

	"time"
)

func (o *Object) run(world *World, timeDelta time.Duration, timeId int64) {
	o.Lock()
	defer o.Unlock()
	defer o.wal.Commit(timeId)

	rotation := o.AngleSpeed * timeDelta.Seconds()

	if rotation != 0 {
		o.Angle = physics.NormalizeAngle(o.Angle + rotation)
		o.Direction = physics.MakeNormalVectorByAngle(o.Angle)
	}

	if o.Speed != 0 {
		o.Velocity = o.Direction.MultiplyOnScalar(o.Speed)
		o.Pos = physics.ApplyVelocityToPosition(&o.Pos, o.Velocity, timeDelta)
		if o.Pos.CheckIfOutOfBounds(0, 0, float64(world.width), float64(world.height)) {
			o.wal.AddDelete()
			return
		}
		o.collisions(world)
	}
	if rotation != 0 || o.Speed != 0 {
		o.wal.AddAngle(o.Angle)
		o.wal.AddRotation(rotation)
		o.wal.AddPosAndVelocityLen(o.Pos, o.Speed)
	}
}

// просчет коллизий с другими объектами
func (o *Object) collisions(w *World) {
	for otherId, otherObject := range w.objects {
		if o.Id == otherId || o.getType() != TypeMissile {
			continue
		}
		if o.isCollideWith(otherObject) {
			o.wal.AddDeleteOtherIds([]uint32{otherId})
			o.wal.AddExplode()
			delete(w.objects, otherId)
			delete(w.objects, o.Id)
			break
		}
	}
}
