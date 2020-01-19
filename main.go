package main

import (
	"aakimov/marsgame/server"
	"aakimov/marsgame/world"

	"log"
	"net/http"
)

func main() {
	s := server.NewServer()
	w := world.NewWorld(s)
	s.Setup()
	go w.Run()
	log.Fatal(http.ListenAndServe(":80", nil))
}
