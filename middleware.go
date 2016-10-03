package broadcaster

type HandlerFunc func(c []*Client, msgs []*Message)

func (m HandlerFunc) ServeHTTP(c []*Client, msgs []*Message) {
	m(c, msgs)
}

type Middleware func(fn HandlerFunc) HandlerFunc
