package broadcaster

import (
	"io"

	"golang.org/x/net/websocket"
)

type Client struct {
	ID    string
	ws    *websocket.Conn
	sv    *Server
	write chan []*Message
	done  chan bool
}

func NewClient(ID string, ws *websocket.Conn, sv *Server, writeBufferSize int) *Client {
	return &Client{
		ID:    ID,
		ws:    ws,
		sv:    sv,
		write: make(chan []*Message, writeBufferSize),
		done:  make(chan bool),
	}
}

func (m *Client) OnWrite(msgs []*Message) {
	m.write <- msgs
}

func (m *Client) Listen() {
	go m.listenWrite()
	m.listenRead()
}

func (m *Client) listenWrite() {
	for {
		select {
		case msgs := <-m.write:
			websocket.JSON.Send(m.ws, msgs)
		case <-m.done:
			m.Done()
			return
		}
	}
}

func (m *Client) listenRead() {
	for {
		select {
		case <-m.done:
			m.Done()
			return
		default:
			msg := []*Message{}
			err := websocket.JSON.Receive(m.ws, &msg)
			if err == io.EOF {
				m.done <- true
			} else if err != nil {
				m.sv.OnError(err)
			} else {
				m.sv.OnBloadCast(msg)
			}
		}
	}
}

func (m *Client) Done() {
	m.sv.OnDelClient(m)
	m.done <- true
}
