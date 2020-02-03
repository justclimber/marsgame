package world

import (
	"aakimov/marsgame/helpers"
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

		mech.Object.Angle += mech.rotateThrottle * MaxRotationValue * throttleRegression
		if mech.Object.Angle > 2*math.Pi {
			mech.Object.Angle = mech.Object.Angle - 2*math.Pi
		} else if mech.Object.Angle < 0 {
			mech.Object.Angle = 2*math.Pi + mech.Object.Angle
		}
		newAngle := mech.Object.Angle
		changeByObject.Angle = &newAngle
		mech.Object.Direction = makeNormalVectorByAngle(newAngle)
		mech.Object.Velocity = mech.Object.Direction.multiplyOnScalar(mech.Object.Velocity.len())
	}

	// просчет движения меха по вектору velocity
	velocityLen := mech.Velocity.len
	if mech.throttle != 0 || velocityLen() != 0 {
		energyNeed := int(mech.throttle * mechFullThrottleEnergyPerSec * timeDelta.Seconds())
		throttleRegression := mech.generator.consumeWithPartlyUsage(energyNeed)
		power := mech.throttle * maxPower * throttleRegression

		newPos, newVelocity := calcMovementObject(&mech.Object, power, timeDelta)
		length := newPos.distanceTo(&mech.Object.Pos)
		mech.Object.Pos = *newPos
		mech.Object.Velocity = newVelocity
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
			p.pickupXelon()
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
		Object: Object{
			Type:            TypeMissile,
			Speed:           MissileSpeed,
			Pos:             missilePos,
			Angle:           missileAngle,
			CollisionRadius: 20,
		},
	}
}
