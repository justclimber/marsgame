package world

import (
	"aakimov/marsgame/helpers"
	"aakimov/marsgame/physics"
	"math"
	"time"
)

const nearTimeDelta = 50
const maxPower = 50000
const mechFullThrottleEnergyPerSec = 5000
const mechFullRotateThrottleEnergyPerSec = 2000
const shootEnergy = 4000
const xelonsInOneCrystal = 200

func areTimeIdNearlySameOrGrater(t1, t2 int64) bool {
	return t1 > t2 || helpers.AbsInt64(t1-t2) < nearTimeDelta
}

func (p *Player) run(timeDelta time.Duration) *ChangeByObject {
	mech := p.mech
	changeByObject := ChangeByObject{
		ObjType: TypePlayer,
		ObjId:   p.id,
	}
	mech.Lock()
	defer mech.Unlock()

	// просчет поворота меха
	if mech.rotateThrottle != 0 {
		energyNeed := int(mech.rotateThrottle * mechFullRotateThrottleEnergyPerSec * timeDelta.Seconds())
		throttleRegression := mech.generator.consumeWithPartlyUsage(energyNeed)

		mech.Obj.Angle += mech.rotateThrottle * MaxRotationValue * throttleRegression
		if mech.Obj.Angle > 2*math.Pi {
			mech.Obj.Angle = mech.Obj.Angle - 2*math.Pi
		} else if mech.Obj.Angle < 0 {
			mech.Obj.Angle = 2*math.Pi + mech.Obj.Angle
		}
		newAngle := mech.Obj.Angle
		changeByObject.Angle = &newAngle
		mech.Obj.Direction = physics.MakeNormalVectorByAngle(newAngle)
		mech.Obj.Velocity = mech.Obj.Direction.MultiplyOnScalar(mech.Obj.Velocity.Len())
	}

	// просчет движения меха по вектору velocity
	velocityLen := mech.Velocity.Len
	if mech.throttle != 0 || velocityLen() != 0 {
		energyNeed := int(mech.throttle * mechFullThrottleEnergyPerSec * timeDelta.Seconds())
		throttleRegression := mech.generator.consumeWithPartlyUsage(energyNeed)
		power := mech.throttle * maxPower * throttleRegression

		newPos, newVelocity := physics.CalcMovementObject(&mech.Obj, power, timeDelta)
		length := newPos.DistanceTo(&mech.Obj.Pos)
		mech.Obj.Pos = *newPos
		mech.Obj.Velocity = newVelocity
		changeByObject.Pos = newPos
		changeByObject.length = &length

		p.collisions(&changeByObject)
	}

	// просчет поворота башни меха
	if mech.cannon.rotateThrottle != 0 {
		mech.cannon.angle += mech.cannon.rotateThrottle * MaxCannonRotationValue
		newCannonAngle := mech.cannon.angle
		changeByObject.CannonAngle = &newCannonAngle
	}

	// просчет выстрела
	if mech.cannon.shoot.state == WillShoot {
		mech.cannon.shoot.state = Planned
		mech.cannon.shoot.willShootAt = p.world.timeId + int64(mech.cannon.shoot.delay)
	}
	if mech.cannon.shoot.state == Planned && areTimeIdNearlySameOrGrater(p.world.timeId, mech.cannon.shoot.willShootAt) {
		if mech.generator.consumeIfHas(shootEnergy) {
			mech.cannon.shoot.state = None
			p.shoot()
		}
	}

	if mech.rotateThrottle != 0 || velocityLen() != 0 || mech.cannon.rotateThrottle != 0 {
		return &changeByObject
	}
	return nil
}

// просчет коллизий с другими объектами
func (p *Player) collisions(ch *ChangeByObject) {
	for id, object := range p.world.objects {
		if object.getType() != TypeMissile && p.mech.isCollideWith(object) {
			ch.DeleteOtherId = id
			delete(p.world.objects, id)
			if object.getType() == TypeXelon {
				p.pickupXelon()
			}
			break
		}
	}
}

// подбор кристалла кселона - увеличение количества кселона в генераторе
func (p *Player) pickupXelon() {
	p.mech.generator.increaseXelons(xelonsInOneCrystal)
}

func (p *Player) shoot() {
	missileAngle := p.mech.cannon.angle + p.mech.Angle
	missilePos := p.mech.Pos

	//move missile a bit of forward far away from mech center
	missilePos.MoveForward(missileAngle, 100.)
	p.world.newObjectsCh <- &Missile{
		Object: NewObject(0,
			TypeMissile,
			missilePos,
			20,
			missileAngle,
			MissileSpeed,
			0,
			10,
		)}
}
