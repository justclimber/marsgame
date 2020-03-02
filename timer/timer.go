package timer

import "time"

type Timer struct {
	valueInSeconds time.Duration
	left           time.Duration
}
