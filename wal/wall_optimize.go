package wal

func (oo *ObjectObserver) optimize() {
	if oo.timeLog.X == defaultInt &&
		oo.timeLog.Y == defaultInt &&
		oo.timeLog.Angle == defaultFloat &&
		oo.timeLog.CannonAngle == defaultFloat &&
		!oo.timeLog.IsNew &&
		!oo.timeLog.Fire &&
		!oo.timeLog.Delete &&
		len(oo.timeLog.DeleteOtherObjectIds) == 0 &&
		!oo.timeLog.Explode {

		oo.timeLog.skip = true
		return
	}
	lastTimeLog := oo.objectLog.LastTimeLog()
	if lastTimeLog == nil {
		oo.lastVelocityLen = &oo.timeLog.VelocityLen
		oo.lastVelocityUntilTimeId = &oo.timeLog.VelocityUntilTimeId
		return
	}
	oo.timeLog.skip = true
	if oo.timeLog.VelocityLen != defaultFloat &&
		oo.timeLog.VelocityLen == *oo.lastVelocityLen &&
		oo.timeLog.VelocityRotation == *oo.lastVelocityRotation {

		lastTimeLog.VelocityUntilTimeId = oo.timeLog.TimeId
		oo.timeLog.X = defaultInt
		oo.timeLog.Y = defaultInt
		oo.timeLog.VelocityLen = defaultFloat
	} else {
		oo.lastVelocityLen = &oo.timeLog.VelocityLen
		oo.lastVelocityUntilTimeId = &oo.timeLog.VelocityUntilTimeId
		oo.timeLog.skip = false
	}

	if oo.timeLog.IsNew ||
		oo.timeLog.Fire ||
		oo.timeLog.Delete ||
		len(oo.timeLog.DeleteOtherObjectIds) != 0 ||
		oo.timeLog.Explode {
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
}
