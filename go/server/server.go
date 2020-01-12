package server

import (
	"log"
)

type Server struct {
	clients         map[string]*Client
	connectClientCh chan *Client
	leaveClientCh   chan *Client
	doneCh          chan bool
	errCh           chan error
	NewClientCh     chan *Client
	SaveAstCodeCh   chan *SaveAstCode
}

func NewServer() *Server {
	return &Server{
		clients:         make(map[string]*Client),
		connectClientCh: make(chan *Client),
		leaveClientCh:   make(chan *Client),
		NewClientCh:     make(chan *Client),
		doneCh:          make(chan bool),
		errCh:           make(chan error),
		SaveAstCodeCh:   make(chan *SaveAstCode),
	}
}

func (s *Server) ListenClients() {
	log.Println("Server start listening")
	c := <-s.connectClientCh
	log.Printf("Client [%s] connected and registered!", c.Id)
	s.clients[c.Id] = c

	for {
		select {
		case c := <-s.connectClientCh:
			log.Printf("Client [%s] registered!", c.Id)
			s.clients[c.Id] = c

		case c := <-s.leaveClientCh:
			log.Printf("Client [%s] unconnected", c.Id)
			delete(s.clients, c.Id)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}

func (s *Server) connectClient(client *Client) {
	s.connectClientCh <- client
	s.NewClientCh <- client
}

type SaveAstCode struct {
	UserId     string
	SourceCode string
}

func (s *Server) saveSourceCode(clientId, sourceCode string) {
	_, ok := s.clients[clientId]
	if !ok {
		log.Fatalf("Save code attempt for inexistant client [%s]", clientId)
	}
	s.SaveAstCodeCh <- &SaveAstCode{
		UserId:     clientId,
		SourceCode: sourceCode,
	}
}

func (s *Server) HandleCommand(client *Client, command *Command) {
	log.Println(command.ToSting())
}
