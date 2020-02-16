package server

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	Id       int
	ws       *websocket.Conn
	server   *Server
	commands chan *Command
	doneCh   chan bool
}

func NewClient(id int, ws *websocket.Conn, server *Server) *Client {
	return &Client{
		Id:       id,
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

func (c *Client) PackAndSendCommand(name string, payload interface{}) {
	c.SendCommand(PackStructToCommand(name, payload))
}

func (c *Client) SendCommand(command *Command) {
	c.commands <- command
}

func (c *Client) listenWrite() {
	for {
		select {
		case cmd := <-c.commands:
			if err := c.ws.WriteJSON(cmd); err != nil {
				log.Println(err)
			}
			//log.Printf("Command sent to the client with payload %s\n", cmd.Payload)
			//log.Printf("Payload size %d\n", len(cmd.Payload))

		case <-c.doneCh:
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

func (c *Client) listenRead() {
	defer func() {
		if err := c.ws.Close(); err != nil {
			log.Println(err)
		}
		c.server.leaveClientCh <- c
	}()
	for {
		select {
		case <-c.doneCh:
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var command Command
			if err := c.ws.ReadJSON(&command); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Reading from ws error: %s\n", err)
				}
				return
			}
			c.server.HandleCommand(c, &command)
		}
	}
}
