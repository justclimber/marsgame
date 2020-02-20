package wal

func (om *ObjectManager) optimize() {
	lastTimeLog := om.objectLog.LastTimeLog()
	if lastTimeLog == nil {
		om.lastVelocityX = om.timeLog.VelocityX
		om.lastVelocityY = om.timeLog.VelocityY
		om.lastVelocityUntilTimeId = om.timeLog.VelocityUntilTimeId
		return
	}
	om.timeLog.skip = true
	if (om.timeLog.VelocityX == nil &&
		om.lastVelocityX == nil &&
		om.lastVelocityRotation == nil) ||
		(om.timeLog.VelocityX != nil &&
			om.lastVelocityX != nil &&
			*om.timeLog.VelocityX == *om.lastVelocityX &&
			*om.timeLog.VelocityY == *om.lastVelocityY) {
		lastTimeLog.VelocityUntilTimeId = &om.timeLog.TimeId
		om.timeLog.X = nil
		om.timeLog.Y = nil
		om.timeLog.VelocityX = nil
		om.timeLog.VelocityY = nil
		om.timeLog.VelocityLen = nil
	} else {
		om.lastVelocityX = om.timeLog.VelocityX
		om.lastVelocityY = om.timeLog.VelocityY
		om.lastVelocityUntilTimeId = om.timeLog.VelocityUntilTimeId
		om.timeLog.skip = false
	}

	//if (om.timeLog.CannonAngle == nil && om.lastVelocityRotation == nil) ||
	//	(om.timeLog.VelocityRotation != nil &&
	//		om.lastVelocityRotation != nil &&
	//		*om.timeLog.VelocityRotation == *om.lastVelocityRotation) {
	//	om.timeLog.Angle = nil
	//	om.timeLog.VelocityRotation = nil
	//} else {
	//	om.timeLog.skip = false
	//}

	//if (om.timeLog.VelocityRotation == nil && om.lastVelocityRotation == nil) ||
	//	(om.timeLog.VelocityRotation != nil &&
	//		om.lastVelocityRotation != nil &&
	//		*om.timeLog.VelocityRotation == *om.lastVelocityRotation) {
	//	om.timeLog.Angle = nil
	//	om.timeLog.VelocityRotation = nil
	//} else {
	//	om.timeLog.skip = false
	//}
}
