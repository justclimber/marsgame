package wal

import (
	"github.com/justclimber/marsgame/server"
)

type Subscriber struct {
	client     *server.Client
	currTimeId int64
}

type Sender struct {
	logStorage    *Log
	terminateCh   chan bool
	logCh         chan *Log
	subscribeCh   chan *server.Client
	unsubscribeCh chan *server.Client
	subscribers   map[uint32]*Subscriber
}

func NewSender(storage *Log) *Sender {
	return &Sender{
		logStorage:    storage,
		terminateCh:   make(chan bool, 1),
		logCh:         make(chan *Log, 10),
		subscribeCh:   make(chan *server.Client, 10),
		unsubscribeCh: make(chan *server.Client, 10),
		subscribers:   make(map[uint32]*Subscriber),
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
			s.subscribers[client.Id] = &Subscriber{client, s.logStorage.currTimeId}
			client.SendBuffer(logToBuffer(s.logStorage))
		case client := <-s.unsubscribeCh:
			delete(s.subscribers, client.Id)
		case log := <-s.logCh:
			s.logStorage.merge(log)
			buf := logToBuffer(log)
			for _, subscriber := range s.subscribers {
				subscriber.client.SendBuffer(buf)
			}
		}
	}
}
