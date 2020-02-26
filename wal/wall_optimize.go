package wal

func (oo *ObjectObserver) optimize() {
	if oo.timeLog.X == defaultInt &&
		oo.timeLog.Y == defaultInt &&
		oo.timeLog.Angle == defaultFloat &&
		oo.timeLog.CannonAngle == defaultFloat &&
		!oo.timeLog.IsNew &&
		!oo.timeLog.Fire &&
		!oo.timeLog.Delete &&
		!oo.timeLog.Explode {

		oo.timeLog.skip = true
		return
	}
	lastTimeLog := oo.objectLog.LastTimeLog()
	if lastTimeLog == nil {
		oo.lastVelocityX = &oo.timeLog.VelocityX
		oo.lastVelocityY = &oo.timeLog.VelocityY
		oo.lastVelocityUntilTimeId = &oo.timeLog.VelocityUntilTimeId
		return
	}
	oo.timeLog.skip = true
	if oo.timeLog.VelocityX != defaultFloat &&
		oo.lastVelocityX != nil &&
		oo.timeLog.VelocityX == *oo.lastVelocityX &&
		oo.timeLog.VelocityY == *oo.lastVelocityY {

		lastTimeLog.VelocityUntilTimeId = oo.timeLog.TimeId
		oo.timeLog.X = defaultInt
		oo.timeLog.Y = defaultInt
		oo.timeLog.VelocityX = defaultFloat
		oo.timeLog.VelocityY = defaultFloat
	} else {
		oo.lastVelocityX = &oo.timeLog.VelocityX
		oo.lastVelocityY = &oo.timeLog.VelocityY
		oo.lastVelocityUntilTimeId = &oo.timeLog.VelocityUntilTimeId
		oo.timeLog.skip = false
	}

	//if (oo.timeLog.CannonAngle == nil && oo.lastVelocityRotation == nil) ||
	//	(oo.timeLog.VelocityRotation != nil &&
	//		oo.lastVelocityRotation != nil &&
	//		*oo.timeLog.VelocityRotation == *oo.lastVelocityRotation) {
	//	oo.timeLog.Angle = nil
	//	oo.timeLog.VelocityRotation = nil
	//} else {
	//	oo.timeLog.skip = false
	//}

	//if (oo.timeLog.VelocityRotation == nil && oo.lastVelocityRotation == nil) ||
	//	(oo.timeLog.VelocityRotation != nil &&
	//		oo.lastVelocityRotation != nil &&
	//		*oo.timeLog.VelocityRotation == *oo.lastVelocityRotation) {
	//	oo.timeLog.Angle = nil
	//	oo.timeLog.VelocityRotation = nil
	//} else {
	//	oo.timeLog.skip = false
	//}
}
