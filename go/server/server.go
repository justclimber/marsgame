package server

import "log"

type Server struct {
	clients         map[string]*Client
	connectClientCh chan *Client
	leaveClientCh   chan *Client
	doneCh          chan bool
	errCh           chan error
}

func NewServer() Server {
	return Server{
		clients:         make(map[string]*Client),
		connectClientCh: make(chan *Client),
		leaveClientCh:   make(chan *Client),
		doneCh:          make(chan bool),
		errCh:           make(chan error),
	}
}

func (s *Server) Listen() {
	log.Println("Server start listening")
	c := <-s.connectClientCh
	log.Printf("Client [%s] connected and registered!", c.id)
	s.clients[c.id] = c

	for {
		select {
		case c := <-s.connectClientCh:
			log.Printf("Client [%s] registered!", c.id)
			s.clients[c.id] = c

		case c := <-s.leaveClientCh:
			delete(s.clients, c.id)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}

func (s *Server) connectClient(client *Client) {
	s.connectClientCh <- client
}

func (s *Server) HandleCommand(client *Client, command *Command) {
	log.Println(command.ToSting())
	client.commands <- command
}
