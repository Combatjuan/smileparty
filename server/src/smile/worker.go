package smile

import (
	"fmt"
	"io"
	"log"

	"golang.org/x/net/websocket"
)

const channelBufSize = 100

var maxId int = 0

type Worker struct {
	id     int
	ws     *websocket.Conn
	server *Server
	ch     chan *SmileLocation
	doneCh chan bool
}

func NewWorker(ws *websocket.Conn, server *Server) *Worker {
	if ws == nil {
		panic("ws cannot be nil")
	}
	if server == nil {
		panic("server cannot be nil")
	}
	maxId++
	ch := make(chan *SmileLocation, channelBufSize)
	doneCh := make(chan bool)

	return &Worker{maxId, ws, server, ch, doneCh}
}

func (c *Worker) Conn() *websocket.Conn {
	return c.ws
}

func (c *Worker) Write(msg *SmileLocation) {
	select {
	case c.ch <- msg:
	default:
		c.server.Del(c)
		err := fmt.Errorf("client $d is disconnected.", c.id)
		c.server.Err(err)
	}
}

func (c *Worker) Done() {
	c.doneCh <- true
}

func (c *Worker) Listen() {
	go c.listenWrite()
	c.listenRead()
}

func (c *Worker) listenWrite() {
	log.Println("Listening write to client")
	for {
		select {
		case msg := <-c.ch:
			log.Println("Send:", msg)
			websocket.JSON.Send(c.ws, msg)
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true
			return
		}
	}
}

func (c *Worker) listenRead() {
	log.Println("Listening read from client")
	for {
		select {
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true
			return
		default:
			var msg SmileLocation
			err := websocket.JSON.Receive(c.ws, &msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				c.server.SendAll(&msg)
			}
		}
	}
}
