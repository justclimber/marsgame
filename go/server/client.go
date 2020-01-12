package server

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	id       string
	ws       *websocket.Conn
	server   *Server
	commands chan *Command
	doneCh   chan bool
}

func NewClient(id string, ws *websocket.Conn, server *Server) Client {
	return Client{
		id:       id,
		ws:       ws,
		server:   server,
		commands: make(chan *Command),
		doneCh:   make(chan bool),
	}
}

func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

func (c *Client) listenWrite() {
	for {
		select {
		case cmd := <-c.commands:
			if err := c.ws.WriteJSON(cmd); err != nil {
				log.Println(err)
			}

		case <-c.doneCh:
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

func (c *Client) listenRead() {
	for {
		select {
		case <-c.doneCh:
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var command Command
			if err := c.ws.ReadJSON(&command); err != nil {
				log.Printf("Unmarshaling error: %s\n", err)
			}
			c.server.HandleCommand(c, &command)
		}
	}
}
