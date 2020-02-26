package server

import (
	"log"
	"strconv"
)

type Server struct {
	clients         map[uint32]*Client
	connectClientCh chan *Client
	leaveClientCh   chan *Client
	doneCh          chan bool
	errCh           chan error
	NewClientCh     chan *Client
	SaveAstCodeCh   chan *SaveAstCode
	ProgramFlowCh   chan *ProgramFlow
	CommandsCh      chan *CommandFromClient
}

func NewServer() *Server {
	return &Server{
		clients:         make(map[uint32]*Client),
		connectClientCh: make(chan *Client, 10),
		leaveClientCh:   make(chan *Client, 10),
		NewClientCh:     make(chan *Client, 10),
		doneCh:          make(chan bool),
		errCh:           make(chan error),
		SaveAstCodeCh:   make(chan *SaveAstCode, 10),
		ProgramFlowCh:   make(chan *ProgramFlow, 10),
		CommandsCh:      make(chan *CommandFromClient, 10),
	}
}

func (s *Server) ListenClients() {
	log.Println("Server start listening")

	for {
		select {
		case c := <-s.connectClientCh:
			log.Printf("Client [%d] registered!", c.Id)
			s.clients[c.Id] = c

		case c := <-s.leaveClientCh:
			log.Printf("Client [%d] unconnected", c.Id)
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
	UserId     uint32
	SourceCode string
}

type ProgramFlowType int

const (
	StopProgram ProgramFlowType = iota
	StartProgram
)

type ProgramFlow struct {
	UserId  uint32          `json:"userId"`
	FlowCmd ProgramFlowType `json:"flowCmd"`
}

type CommandFromClient struct {
	UserId  uint32
	Command *Command
}

func (s *Server) HandleCommand(client *Client, command *Command) {
	switch command.Type {
	case "saveCode":
		s.SaveAstCodeCh <- &SaveAstCode{
			UserId:     client.Id,
			SourceCode: command.Payload,
		}
	case "programFlow":
		flowAsInt, _ := strconv.Atoi(command.Payload)
		s.ProgramFlowCh <- &ProgramFlow{
			UserId:  client.Id,
			FlowCmd: ProgramFlowType(flowAsInt),
		}
	default:
		s.CommandsCh <- &CommandFromClient{client.Id, command}
	}
}
