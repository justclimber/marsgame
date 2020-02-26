package wal

import "aakimov/marsgame/server"

type Sender struct {
	terminateCh   chan bool
	logCh         chan *Log
	subscribeCh   chan *server.Client
	unsubscribeCh chan *server.Client
	clients       map[int]*server.Client
}

func NewSender() *Sender {
	return &Sender{
		terminateCh:   make(chan bool, 1),
		logCh:         make(chan *Log, 10),
		subscribeCh:   make(chan *server.Client, 10),
		unsubscribeCh: make(chan *server.Client, 10),
		clients:       make(map[int]*server.Client),
	}
}

func (s *Sender) Subscribe(client *server.Client) {
	s.subscribeCh <- client
}

func (s *Sender) Unsubscribe(client *server.Client) {
	s.unsubscribeCh <- client
}

func (s *Sender) Terminate() {
	s.terminateCh <- true
}

func (s *Sender) SendLoop() {
	for {
		select {
		case <-s.terminateCh:
			return
		case client := <-s.subscribeCh:
			s.clients[client.Id] = client
		case client := <-s.unsubscribeCh:
			delete(s.clients, client.Id)
		case log := <-s.logCh:
			for _, client := range s.clients {
				command := server.PackStructToCommand("wal", log)
				client.SendCommand(command)
			}
		}
	}
}
