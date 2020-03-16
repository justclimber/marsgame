package timer

import (
	"github.com/justclimber/marsgame/flatbuffers/generated/InitBuffers"
	"log"

	"time"
)

type Timer struct {
	value     time.Duration
	realValue time.Duration
	stopCh    chan bool
	state     int8
	timer     *time.Timer
	startedAt time.Time
	pausedAt  time.Time
}

func (t *Timer) Value() time.Duration {
	return t.value
}

func NewTimer(value time.Duration, timeMultiplicator int) *Timer {
	return &Timer{
		value:     value,
		realValue: value / time.Duration(timeMultiplicator),
		stopCh:    make(chan bool),
	}
}

func (t *Timer) Start(onStopCallback func()) {
	if t.state != InitBuffers.TimerStateStopped {
		log.Fatal("Attempt to resume not stopped timer")
	}
	t.startedAt = time.Now()
	t.state = InitBuffers.TimerStateStarted
	t.timer = time.AfterFunc(t.realValue, func() {
		onStopCallback()
		log.Println("Timer end!")
		t.state = InitBuffers.TimerStateExpired
	})
}

func (t *Timer) IsStopped() bool {
	return t.state == InitBuffers.TimerStateStopped
}

func (t *Timer) Stop() {
	if t.state != InitBuffers.TimerStateStarted {
		log.Fatal("Attempt to stop not started timer")
	}
	t.timer.Stop()
	t.state = InitBuffers.TimerStateStopped
}

func (t *Timer) Pause() {
	if t.state != InitBuffers.TimerStateStarted {
		log.Fatal("Attempt to resume not started timer")
	}
	t.timer.Stop()
	t.state = InitBuffers.TimerStatePaused
	t.pausedAt = time.Now()
}

func (t *Timer) Resume() {
	if t.state != InitBuffers.TimerStatePaused {
		log.Fatal("Attempt to resume not paused timer")
	}
	t.state = InitBuffers.TimerStatePaused
	t.timer.Reset(t.realValue - t.pausedAt.Sub(t.startedAt))
}

func (t *Timer) State() int8 {
	return t.state
}
