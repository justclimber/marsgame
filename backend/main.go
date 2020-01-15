package main

import (
	"aakimov/marsgame/backend/server"
	"aakimov/marsgame/backend/world"

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
