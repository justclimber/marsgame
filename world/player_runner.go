package world

import (
	"aakimov/marsgame/helpers"
	"aakimov/marsgame/physics"
	"time"
)

const nearTimeDelta = 50
const maxPower = 50000
const mechFullThrottleEnergyPerSec = 5000
const mechFullRotateThrottleEnergyPerSec = 2000
const shootEnergy = 4000
const xelonsInOneCrystal = 200
const MissileSpeed = 50
const MaxRotationValue float64 = 0.1
const MaxCannonRotationValue float64 = 0.11

func areTimeIdNearlySameOrGrater(t1, t2 int64) bool {
	return t1 > t2 || helpers.AbsInt64(t1-t2) < nearTimeDelta
}

func (p *Player) run(timeDelta time.Duration, timeId int64) {
	mech := p.mech
	mech.Lock()
	defer mech.Unlock()
	defer p.wal.Commit(timeId)

	// просчет поворота меха
	if mech.rotateThrottle != 0 {
		energyNeed := int(mech.rotateThrottle * mechFullRotateThrottleEnergyPerSec * timeDelta.Seconds())
		throttleRegression := mech.generator.consumeWithPartlyUsage(energyNeed)

		rotation := mech.rotateThrottle * MaxRotationValue * throttleRegression
		mech.Obj.Angle = physics.NormalizeAngle(mech.Obj.Angle + rotation)
		mech.Obj.Direction = physics.MakeNormalVectorByAngle(mech.Obj.Angle)
		mech.Obj.Velocity = mech.Obj.Direction.MultiplyOnScalar(mech.Obj.Velocity.Len())

		p.wal.AddRotation(rotation)
		p.wal.AddAngle(mech.Obj.Angle)
	}

	// просчет движения меха по вектору velocity
	velocityLen := mech.Velocity.Len
	if mech.throttle != 0 || velocityLen() != 0 {
		energyNeed := int(mech.throttle * mechFullThrottleEnergyPerSec * timeDelta.Seconds())
		throttleRegression := mech.generator.consumeWithPartlyUsage(energyNeed)
		power := mech.throttle * maxPower * throttleRegression

		mech.Obj.Pos, mech.Obj.Velocity = physics.MoveObjectByForces(&mech.Obj, power, timeDelta)
		p.wal.AddPosAndVelocity(mech.Obj.Pos, mech.Obj.Velocity)

		p.collisions()
	}

	// просчет поворота башни меха
	if mech.cannon.rotateThrottle != 0 {
		cannonRotation := mech.cannon.rotateThrottle * MaxCannonRotationValue
		mech.cannon.angle += cannonRotation
		p.wal.AddCannonAngle(cannonRotation)
		p.wal.AddCannonAngle(mech.cannon.angle)
	}

	// просчет выстрела
	if mech.cannon.shoot.state == WillShoot {
		mech.cannon.shoot.state = Planned
		mech.cannon.shoot.willShootAt = timeId + int64(mech.cannon.shoot.delay)
	}
	if mech.cannon.shoot.state == Planned && areTimeIdNearlySameOrGrater(timeId, mech.cannon.shoot.willShootAt) {
		if mech.generator.consumeIfHas(shootEnergy) {
			mech.cannon.shoot.state = None
			p.shoot()
			p.wal.AddShoot()
		}
	}
}

// просчет коллизий с другими объектами
func (p *Player) collisions() {
	for id, object := range p.world.objects {
		if object.getType() != TypeMissile && p.mech.isCollideWith(object) {
			p.wal.AddDeleteOther(id)
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
