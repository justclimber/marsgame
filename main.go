package main

import (
	"github.com/justclimber/marsgame/server"
	"github.com/justclimber/marsgame/world"

	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	s := server.NewServer()
	w := world.NewWorld(s)
	w.Bootstrap()
	s.Setup()
}
