package world

import (
	"log"
	"sync"
	"time"
)

type Generator struct {
	sync.Mutex
	value      int
	maxValue   int
	xelons     int
	efficiency int
	rateMs     time.Duration
}

func (g *Generator) geValue() int {
	g.Lock()
	defer g.Unlock()
	return g.value
}

// если запрашиваемой энергии не хватает - не потребляем и возвращаем false
// в противном случае - потребляем и возвращаем true
func (g *Generator) consumeIfHas(value int) bool {
	g.Lock()
	defer g.Unlock()

	if g.value-value < 0 {
		return false
	}
	g.value -= value
	return true
}

// если запрашиваемой энергии не хватает - возвращаем в процентах, на сколько ее хватило
func (g *Generator) consumeWithPartlyUsage(value int) float64 {
	g.Lock()
	defer g.Unlock()

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
	g.Lock()
	g.value -= value
	if g.value <= 0 {
		needWait = true
	}
	g.Unlock()

	if needWait {
		ticker := time.NewTicker(g.rateMs * time.Millisecond)
		defer ticker.Stop()
		stop := make(chan bool, 1)
		for {
			select {
			case <-ticker.C:
				g.Lock()
				if g.value > 0 {
					log.Println("cpu trhottling", g.value)
					stop <- true
				}
				g.Unlock()
			case <-stop:
				return
			}
		}
	}
}

func (g *Generator) increaseXelons(incr int) {
	g.Lock()
	defer g.Unlock()
	g.xelons += incr
}

func (g *Generator) consumeXelons() int {
	if g.xelons > 0 {
		g.xelons -= 1
	}
	return 1
}

func (g *Generator) start() {
	ticker := time.NewTicker(g.rateMs * time.Millisecond)
	for range ticker.C {
		g.Lock()
		if g.value > g.maxValue {
			g.value = g.maxValue
		} else {
			g.value += g.consumeXelons() * g.efficiency
		}
		log.Println(g.xelons)
		g.Unlock()
	}
}
