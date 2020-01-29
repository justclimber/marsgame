package world

import "time"

// силя тяготения
const G = 5.

// коэффициент трения
const CoeffFriction = 1.2

// коэффициент сопротивления воздуха
const CoeffAirResist = 1.5

// рассчет силы тяги
func calcTractionForce(direction Vector, enginePower float64) Vector {
	return direction.multiplyOnScalar(enginePower)
}

// рассчет силы сопротивления воздуха
func calcAirResistForce(velocity Vector) Vector {
	return velocity.multiplyOnScalar(-CoeffAirResist * velocity.len())
}

// расчет силы трения
func calcFrictionForce(direction Vector, m float64) Vector {
	return direction.multiplyOnScalar(-CoeffFriction * m * G)
}

// расчет ускорения
func calcAccelerate(force Vector, m float64) Vector {
	return force.multiplyOnScalar(1 / m)
}

// рассчет скорости
func applyAccelerateToVelocity(v Vector, a Vector, dt time.Duration) Vector {
	return v.add(a.multiplyOnScalar(dt.Seconds()))
}

// рассчет перемещения
func applyVelocityToPosition(p Point, v Vector, dt time.Duration) Point {
	return p.add(v.multiplyOnScalar(dt.Seconds()))
}

// общий рассчет
func calculateMovement(p Point, v Vector, d Vector, power, m float64, dt time.Duration) Point {
	tractionForce := calcTractionForce(d, power)
	airResistForce := calcAirResistForce(v)
	frictionForce := calcFrictionForce(d, m)
	force := tractionForce.add(airResistForce)
	force = force.add(frictionForce)

	accelerate := calcAccelerate(force, m)
	vNew := applyAccelerateToVelocity(v, accelerate, dt)

	// сила трения на малых скоростях может привести к отризательной скорости, убираем это
	if v.multiplyOnVector(vNew) < 0 {
		return p
	}

	return applyVelocityToPosition(p, vNew, dt)
}
