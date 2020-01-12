package main

import (
	"aakimov/marsgame/go/server"
	"math/rand"
	"time"

	"log"
	"net/http"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	log.Println("Running a http server")
	server.SetupRoutes()
	log.Fatal(http.ListenAndServe(":80", nil))
}
