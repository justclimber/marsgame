package world

import (
	"aakimov/marsgame/helpers"
	"math"
)

const nearTimeDelta = 50

func areTimeIdNearlySameOrGrater(t1, t2 int64) bool {
	return t1 > t2 || helpers.AbsInt64(t1-t2) < nearTimeDelta
}

func (p *Player) run(world *World) *ChangeByObject {
	mech := p.mech
	changeByObject := ChangeByObject{
		ObjType: TypePlayer,
		ObjId:   p.id,
	}
	mech.Lock()
	defer mech.Unlock()

	if mech.cannon.shoot.state == WillShoot {
		mech.cannon.shoot.state = Planned
		mech.cannon.shoot.willShootAt = world.timeId + int64(mech.cannon.shoot.delay)
	}
	if mech.cannon.shoot.state == Planned && areTimeIdNearlySameOrGrater(world.timeId, mech.cannon.shoot.willShootAt) {
		mech.cannon.shoot.state = None
		p.shoot(world)
	}
	if mech.rotateThrottle != 0 {
		mech.Object.Angle += mech.rotateThrottle * MaxRotationValue
		if mech.Object.Angle > 2*math.Pi {
			mech.Object.Angle = mech.Object.Angle - 2*math.Pi
		} else if mech.Object.Angle < 0 {
			mech.Object.Angle = 2*math.Pi + mech.Object.Angle
		}
		newAngle := mech.Object.Angle
		changeByObject.Angle = &newAngle
	}
	if mech.throttle != 0 {
		length := mech.throttle * MaxMovingLength
		mech.Object.Pos.MoveForward(mech.Object.Angle, length)
		newPos := mech.Object.Pos
		changeByObject.Pos = &newPos
		changeByObject.length = &length
	}
	if mech.cannon.rotateThrottle != 0 {
		mech.cannon.angle += mech.cannon.rotateThrottle * MaxCannonRotationValue
		newCannonAngle := mech.cannon.angle
		changeByObject.CannonAngle = &newCannonAngle
	}

	if mech.rotateThrottle != 0 || mech.throttle != 0 || mech.cannon.rotateThrottle != 0 {
		return &changeByObject
	}
	return nil
}

func (p *Player) shoot(world *World) {
	missileAngle := p.mech.cannon.angle + p.mech.Angle
	missilePos := p.mech.Pos

	//move missile a bit of forward far away from mech center
	missilePos.MoveForward(missileAngle, 100.)
	world.newObjectsCh <- &Missile{
		Object: Object{
			Type:            TypeMissile,
			Speed:           MissileSpeed,
			Pos:             missilePos,
			Angle:           missileAngle,
			CollisionRadius: 20,
		},
	}
}
