package main

import (
	"aakimov/marsgame/server"
	"aakimov/marsgame/world"
)

func main() {
	s := server.NewServer()
	w := world.NewWorld(s)
	go w.Run()
	s.Setup()
}
