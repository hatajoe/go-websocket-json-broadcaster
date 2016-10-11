package broadcaster

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	addClient chan *Client
	delClient chan *Client
	err       chan error
	broadcast chan []*Message
	messages  []*Message
}

func NewServer() *Server {
	return &Server{
		addClient: make(chan *Client),
		delClient: make(chan *Client),
		err:       make(chan error),
		broadcast: make(chan []*Message),
		messages:  []*Message{},
	}
}

func (m *Server) OnAddClient(c *Client) {
	m.addClient <- c
}

func (m *Server) OnDelClient(c *Client) {
	m.delClient <- c
}

func (m *Server) OnError(err error) {
	m.err <- err
}

func (m *Server) OnBloadCast(msgs []*Message) {
	m.broadcast <- msgs
}

func (m *Server) Listen(path string, writeBufferSize int, mw Middleware) {
	http.Handle(path, websocket.Handler(func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				m.err <- err
			}
		}()
		c := NewClient(MustUUID(), ws, m, writeBufferSize)
		m.OnAddClient(c)
		c.Listen()
	}))

	clients := []*Client{}
	for {
		select {
		case c := <-m.addClient:
			clients = append(clients, c)
			c.OnWrite(m.messages)
		case c := <-m.delClient:
			for i, v := range clients {
				if v.ID != c.ID {
					continue
				}
				clients = append(clients[:i], clients[i:]...)
				break
			}
		case err := <-m.err:
			log.Println(err)
		case msgs := <-m.broadcast:
			m.messages = append(m.messages, msgs...)
			mw(HandlerFunc(func(c []*Client, ms []*Message) {
				for _, c := range c {
					c.OnWrite(ms)
				}
			}))(clients, msgs)
		default:
		}
	}
}
