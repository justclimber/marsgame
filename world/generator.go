package world

import (
	"log"
	"sync"
	"time"
)

type Generator struct {
	gmu       sync.Mutex
	value     int
	maxValue  int
	increment int
	rateMs    time.Duration
}

func (g *Generator) geValue() int {
	g.gmu.Lock()
	defer g.gmu.Unlock()
	return g.value
}

func (g *Generator) consume(value int) {
	needWait := false
	g.gmu.Lock()
	g.value -= value
	if g.value <= 0 {
		needWait = true
	}
	g.gmu.Unlock()

	if needWait {
		ticker := time.NewTicker(g.rateMs * time.Millisecond)
		stop := make(chan bool, 1)
		for {
			select {
			case <-ticker.C:
				g.gmu.Lock()
				if g.value > 0 {
					log.Println("cpu trhottling", g.value)
					stop <- true
					ticker.Stop()
				}
				g.gmu.Unlock()
			case <-stop:
				return
			}
		}
	}
}

func (g *Generator) start() {
	ticker := time.NewTicker(g.rateMs * time.Millisecond)
	for range ticker.C {
		g.gmu.Lock()
		g.value += g.increment
		if g.value > g.maxValue {
			g.value = g.maxValue
		}
		g.gmu.Unlock()
	}
}
