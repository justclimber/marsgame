package server

import (
	"encoding/json"
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

type saveCode struct {
	UserId     string `json:"userId"`
	SourceCode string `json:"sourceCode"`
}

func (s *Server) saveSourceCodeEndpoint(w http.ResponseWriter, r *http.Request) {
	var sc saveCode
	err := json.NewDecoder(r.Body).Decode(&sc)
	if err != nil {
		log.Fatal(err)
	}

	s.saveSourceCode(sc.UserId, sc.SourceCode)
}

func (s *Server) programFlowCmdEndpoint(w http.ResponseWriter, r *http.Request) {
	var pf ProgramFlow
	err := json.NewDecoder(r.Body).Decode(&pf)
	if err != nil {
		log.Fatal(err)
	}

	s.programFlowCmd(pf.UserId, pf.FlowCmd)
}

func (s *Server) Setup() {
	go s.ListenClients()
	http.Handle("/", http.FileServer(http.Dir("./frontend/static")))
	http.HandleFunc("/ws", s.wsEndpoint)
	http.HandleFunc("/save_source_code", s.saveSourceCodeEndpoint)
	http.HandleFunc("/program_flow", s.programFlowCmdEndpoint)
}
