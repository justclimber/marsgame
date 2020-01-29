package world

import "time"

// силя тяготения
const G = 5.

// коэффициент трения
const CoeffFriction = 1

// коэффициент сопротивления воздуха
const CoeffAirResist = 1

// рассчет силы тяги
func calcTractionForce(direction *Vector, enginePower float64) *Vector {
	return direction.multiplyOnScalar(enginePower)
}

// рассчет силы сопротивления воздуха
func calcAirResistForce(velocity *Vector) *Vector {
	return velocity.multiplyOnScalar(-CoeffAirResist * velocity.len())
}

// расчет силы трения
func calcFrictionForce(direction *Vector, weight float64) *Vector {
	return direction.multiplyOnScalar(-CoeffFriction * weight * G)
}

// расчет ускорения
func calcAccelerate(force *Vector, weight float64) *Vector {
	return force.multiplyOnScalar(1 / weight)
}

// рассчет скорости
func applyAccelerateToVelocity(v *Vector, a *Vector, dt time.Duration) *Vector {
	return v.add(a.multiplyOnScalar(dt.Seconds()))
}

// рассчет перемещения
func applyVelocityToPosition(p *Point, v *Vector, dt time.Duration) *Point {
	return p.add(v.multiplyOnScalar(dt.Seconds()))
}

// общий рассчет
func calcMovementObject(obj *Object, power float64, dt time.Duration) (*Point, *Vector) {
	//m := make(map[string]interface{})
	//m["dir"] = obj.Direction.X
	//m["velocity"] = obj.Velocity.X
	//m["power"] = power
	//m["weight"] = obj.Weight
	//m["dt"] = dt.Seconds()
	//prettyPrint("input", m)
	tractionForce := calcTractionForce(obj.Direction, power)
	airResistForce := calcAirResistForce(obj.Velocity)
	frictionForce := calcFrictionForce(obj.Direction, obj.Weight)
	force := tractionForce.add(airResistForce).add(frictionForce)

	accelerate := calcAccelerate(force, obj.Weight)
	vNew := applyAccelerateToVelocity(obj.Velocity, accelerate, dt)

	m1 := make(map[string]interface{})
	m1["tractionForce"] = tractionForce.X
	m1["airResistForce"] = airResistForce.X
	m1["frictionForce"] = frictionForce.X
	m1["force"] = force.X
	m1["velocity"] = vNew.X
	m1["accelerate"] = accelerate.X
	prettyPrint("output", m1)

	// сила трения на малых скоростях может привести к отризательной скорости, убираем это
	if obj.Direction.multiplyOnVector(vNew) < 0 {
		return &obj.Pos, &Vector{}
	}

	return applyVelocityToPosition(&obj.Pos, vNew, dt), vNew
}
