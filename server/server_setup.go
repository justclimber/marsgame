package server

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *Server) wsEndpoint(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client open ws connection")
	client := NewClient(id, ws, s)
	s.connectClient(client)
	client.Listen()
}

func (s *Server) Setup() {
	go s.ListenClients()
	http.Handle("/", http.FileServer(http.Dir("./frontend/static")))
	http.HandleFunc("/ws", s.wsEndpoint)
}
