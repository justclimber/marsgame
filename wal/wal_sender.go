package wal

import "aakimov/marsgame/server"

type WallSender struct {
	terminateCh   chan bool
	logCh         chan *Log
	subscribeCh   chan *server.Client
	unsubscribeCh chan *server.Client
	clients       map[int]*server.Client
}

func NewWallSender() *WallSender {
	return &WallSender{
		terminateCh:   make(chan bool, 1),
		logCh:         make(chan *Log, 10),
		subscribeCh:   make(chan *server.Client, 10),
		unsubscribeCh: make(chan *server.Client, 10),
		clients:       make(map[int]*server.Client),
	}
}

func (ws *WallSender) Subscribe(client *server.Client) {
	ws.subscribeCh <- client
}

func (ws *WallSender) Unsubscribe(client *server.Client) {
	ws.unsubscribeCh <- client
}

func (ws *WallSender) Terminate() {
	ws.terminateCh <- true
}

func (ws *WallSender) SendLoop() {
	for {
		select {
		case <-ws.terminateCh:
			return
		case client := <-ws.subscribeCh:
			ws.clients[client.Id] = client
		case client := <-ws.unsubscribeCh:
			delete(ws.clients, client.Id)
		case log := <-ws.logCh:
			for _, client := range ws.clients {
				command := server.PackStructToCommand("wal", log)
				client.SendCommand(command)
			}
		}
	}
}
