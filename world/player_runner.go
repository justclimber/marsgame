package world

import (
	"aakimov/marsgame/helpers"
	"aakimov/marsgame/physics"
	"time"
)

const nearTimeDelta = 50
const maxPower = 30000
const mechFullThrottleEnergyPerSec = 5000
const mechFullRotateThrottleEnergyPerSec = 2000
const shootEnergy = 4000
const xelonsInOneCrystal = 200
const MissileSpeed = 280
const MaxRotationValue float64 = 1
const MaxCannonRotationValue float64 = 1.1

func areTimeIdNearlySameOrGrater(t1, t2 int64) bool {
	return t1 > t2 || helpers.AbsInt64(t1-t2) < nearTimeDelta
}

func (p *Player) run(timeDelta time.Duration, timeId int64) {
	p.mech.Lock()
	defer p.mech.Unlock()
	defer p.wal.Commit(timeId)

	mech := p.mech
	commands := mech.commands
	rotation := 0.
	cannonRotation := 0.

	// просчет поворота меха
	if commands.mech.rotate != 0 {
		energyNeed := int(commands.mech.rotate * mechFullRotateThrottleEnergyPerSec * timeDelta.Seconds())
		throttleRegression := mech.generator.consumeWithPartlyUsage(energyNeed)

		rotation = commands.mech.rotate * MaxRotationValue * throttleRegression * timeDelta.Seconds()
		mech.Obj.Angle = physics.NormalizeAngle(mech.Obj.Angle + rotation)
		mech.Obj.Direction = physics.MakeNormalVectorByAngle(mech.Obj.Angle)
		mech.Obj.Velocity = mech.Obj.Direction.MultiplyOnScalar(mech.Obj.Velocity.Len())
	}

	// просчет движения меха по вектору velocity
	velocityLen := mech.Velocity.Len()
	if commands.mech.move != 0 || velocityLen != 0 {
		energyNeed := int(commands.mech.move * mechFullThrottleEnergyPerSec * timeDelta.Seconds())
		throttleRegression := mech.generator.consumeWithPartlyUsage(energyNeed)
		power := commands.mech.move * maxPower * throttleRegression

		mech.Obj.Pos, mech.Obj.Velocity = physics.MoveObjectByForces(&mech.Obj, power, timeDelta)
		p.collisions()
	}

	// просчет поворота башни меха
	if commands.cannon.rotate != 0 {
		cannonRotation = commands.cannon.rotate * MaxCannonRotationValue * timeDelta.Seconds()
		mech.cannon.angle += cannonRotation
	}

	if velocityLen != 0 || rotation != 0 {
		p.wal.AddPosAndVelocityLen(mech.Obj.Pos, velocityLen)
		p.wal.AddRotation(rotation)
		p.wal.AddAngle(mech.Obj.Angle)
		p.wal.AddCannonAngle(mech.cannon.angle)
		p.wal.AddCannonRotation(cannonRotation)
	}

	// просчет выстрела
	if commands.cannon.shoot.state == WillShoot {
		commands.cannon.shoot.state = None
		mech.cannon.shoot.state = Planned
		mech.cannon.shoot.willShootAt = timeId + int64(commands.cannon.shoot.delay)
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
			p.wal.AddDeleteOtherIds([]uint32{id})
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
	missilePos.MoveForward(missileAngle, 20.)
	p.world.newObjectsCh <- &Missile{
		Object: NewObject(0,
			TypeMissile,
			missilePos,
			5,
			missileAngle,
			MissileSpeed,
			0,
			10,
		)}
}
