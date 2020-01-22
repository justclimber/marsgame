package main

import (
	"aakimov/marsgame/server"
	"aakimov/marsgame/world"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	s := server.NewServer()
	w := world.NewWorld(s)
	w.MakeRandomObjects()
	go w.Run()
	s.Setup()
}
