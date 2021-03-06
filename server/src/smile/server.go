package smile

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"time"
)

type Server struct {
	messages    []*SmileLocation
	clients     map[int]*Worker
	addCh       chan *Worker
	delCh       chan *Worker
	sendAllCh   chan *SmileLocation
	doneCh      chan bool
	errCh       chan error
	ticker      *time.Ticker
	currentTime int
}

func NewServer() *Server {
	messages := make([]*SmileLocation, 0)
	clients := make(map[int]*Worker)
	addCh := make(chan *Worker)
	delCh := make(chan *Worker)
	sendAllCh := make(chan *SmileLocation)
	doneCh := make(chan bool)
	errCh := make(chan error)
	ticker := time.NewTicker(time.Second * 10)
	var currentTime int

	return &Server{
		messages, clients, addCh, delCh, sendAllCh, doneCh, errCh, ticker, currentTime,
	}
}

func (s *Server) Add(c *Worker) {
	s.addCh <- c
}

func (s *Server) Del(c *Worker) {
	s.delCh <- c
}

func (s *Server) SendAll(msg *SmileLocation) {
	s.sendAllCh <- msg
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendAll(msg *SmileLocation) {
	for _, c := range s.clients {
		c.Write(msg)
	}
}

func (s *Server) Listen() {
	log.Println("Listening server...")

	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()
		client := NewWorker(ws, s)
		s.Add(client)
		client.Listen()
	}
	http.Handle("/start", websocket.Handler(onConnected))
	log.Println("Created handler")

	for {
		select {
		case c := <-s.addCh:
			log.Println("Added new worker")
			s.clients[c.id] = c
			log.Println("Now", len(s.clients), " connections.")
		case c := <-s.delCh:
			log.Println("Delete worker")
			delete(s.clients, c.id)
		case <-s.ticker.C:
			s.currentTime += 10
			if s.currentTime > 1000 {
				s.currentTime = 0
			}
			s.sendAll(&SmileLocation{"tick", s.currentTime, s.currentTime})
			log.Println("Tick.")
		case msg := <-s.sendAllCh:
			log.Println("Send all:", msg)
			s.messages = append(s.messages, msg)
			s.sendAll(msg)
		case err := <-s.errCh:
			log.Println("Error:", err.Error())
		case <-s.doneCh:
			return
		}
	}
}
