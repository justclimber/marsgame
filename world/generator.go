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

// если запрашиваемой энергии не хватает - не потребляем и возвращаем false
// в противном случае - потребляем и возвращаем true
func (g *Generator) consumeIfHas(value int) bool {
	g.gmu.Lock()
	defer g.gmu.Unlock()

	if g.value-value < 0 {
		return false
	}
	g.value -= value
	return true
}

// если запрашиваемой энергии не хватает - возвращаем в процентах, на сколько ее хватило
func (g *Generator) consumeWithPartlyUsage(value int) float64 {
	g.gmu.Lock()
	defer g.gmu.Unlock()

	g.value -= value
	if g.value <= 0 {
		regression := 1. + float64(g.value)/float64(value)
		g.value = 0
		return regression
	}
	return 1.
}

// если запрашиваемой энергии не хватает - запускаем таймер и ждем пока ее хватит
// и только после того как смогли потребить эту энергию - выходим из метода
func (g *Generator) consumeWithThrottling(value int) {
	needWait := false
	g.gmu.Lock()
	g.value -= value
	if g.value <= 0 {
		needWait = true
	}
	g.gmu.Unlock()

	if needWait {
		ticker := time.NewTicker(g.rateMs * time.Millisecond)
		defer ticker.Stop()
		stop := make(chan bool, 1)
		for {
			select {
			case <-ticker.C:
				g.gmu.Lock()
				if g.value > 0 {
					log.Println("cpu trhottling", g.value)
					stop <- true
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
