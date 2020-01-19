package server

import (
	"encoding/json"
	"fmt"
	engineio "github.com/googollee/go-engine.io"
	"github.com/googollee/go-engine.io/transport"
	"github.com/googollee/go-engine.io/transport/polling"
	"github.com/googollee/go-engine.io/transport/websocket"
	socketio "github.com/googollee/go-socket.io"
	"github.com/rs/cors"
	"log"
	"net/http"
)

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

func getWsServer() *socketio.Server {
	pt := polling.Default
	wt := websocket.Default
	wt.CheckOrigin = func(req *http.Request) bool {
		return true
	}

	wsServer, err := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{pt, wt},
	})
	if err != nil {
		log.Fatal(err)
	}
	return wsServer
}

func (s *Server) Setup() {
	go s.ListenClients()

	mux := http.NewServeMux()

	s.wsServer = getWsServer()

	s.wsServer.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
	s.wsServer.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})
	s.wsServer.OnEvent("/", "command", func(s socketio.Conn, msg string) {
		cmd := Command{
			Type:    "TestCommand",
			Payload: msg,
		}
		fmt.Printf("command received: %s\n", cmd)
	})
	s.wsServer.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})
	s.wsServer.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})
	s.wsServer.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})
	s.wsServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})
	go s.wsServer.Serve()
	defer s.wsServer.Close()

	mux.Handle("/socket.io/", s.wsServer)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedMethods:   []string{"GET", "PUT", "OPTIONS", "POST", "DELETE"},
		AllowCredentials: true,
	})

	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", c.Handler(mux)))

	//http.Handle("/", http.FileServer(http.Dir("./frontend/static")))
	//http.HandleFunc("/ws", s.wsEndpoint)
	//http.HandleFunc("/save_source_code", s.saveSourceCodeEndpoint)
	//http.HandleFunc("/program_flow", s.programFlowCmdEndpoint)
}
