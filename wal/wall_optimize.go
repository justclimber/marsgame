package wal

func (oo *ObjectObserver) optimize() {
	lastTimeLog := oo.objectLog.LastTimeLog()
	if lastTimeLog == nil {
		oo.lastVelocityX = oo.timeLog.VelocityX
		oo.lastVelocityY = oo.timeLog.VelocityY
		oo.lastVelocityUntilTimeId = oo.timeLog.VelocityUntilTimeId
		return
	}
	oo.timeLog.skip = true
	if (oo.timeLog.VelocityX == nil &&
		oo.lastVelocityX == nil &&
		oo.lastVelocityRotation == nil) ||
		(oo.timeLog.VelocityX != nil &&
			oo.lastVelocityX != nil &&
			*oo.timeLog.VelocityX == *oo.lastVelocityX &&
			*oo.timeLog.VelocityY == *oo.lastVelocityY) {
		lastTimeLog.VelocityUntilTimeId = &oo.timeLog.TimeId
		oo.timeLog.X = nil
		oo.timeLog.Y = nil
		oo.timeLog.VelocityX = nil
		oo.timeLog.VelocityY = nil
	} else {
		oo.lastVelocityX = oo.timeLog.VelocityX
		oo.lastVelocityY = oo.timeLog.VelocityY
		oo.lastVelocityUntilTimeId = oo.timeLog.VelocityUntilTimeId
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
