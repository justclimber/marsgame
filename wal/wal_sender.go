package wal

import (
	"aakimov/marsgame/flatbuffers/WalBuffers"
	"aakimov/marsgame/server"
)
import flatbuffers "github.com/google/flatbuffers/go"

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
			buf := s.logToBuffer(log)
			for _, client := range s.clients {
				client.SendBuffer(buf)
			}
		}
	}
}

func (s *Sender) logToBuffer(logToBuff *Log) []byte {
	timelog1 := logToBuff.Objects[2].Times[0]
	builder := flatbuffers.NewBuilder(1024)
	WalBuffers.TimeLogStart(builder)
	if timelog1.X != nil {
		WalBuffers.TimeLogAddX(builder, int32(*timelog1.X))
		WalBuffers.TimeLogAddY(builder, int32(*timelog1.Y))
	}
	if timelog1.Angle != nil {
		WalBuffers.TimeLogAddAngle(builder, float32(*timelog1.Angle))
	}
	timeLogBuffer := WalBuffers.TimeLogEnd(builder)
	builder.Finish(timeLogBuffer)
	buf := builder.FinishedBytes()

	return buf
}
