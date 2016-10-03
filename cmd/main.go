package main

import (
	"flag"
	"fmt"
	"net/http"

	broadcaster "github.com/hatajoe/go-websocket-json-broadcaster"
)

func main() {
	port := flag.Int("p", 9218, "websocket listen port")
	endpoint := flag.String("e", "/", "websocket application endpoint path")
	bfSize := flag.Int("s", 100, "message buffer size per client")
	flag.Parse()

	sv := broadcaster.NewServer()
	go sv.Listen(
		*endpoint,
		*bfSize,
		broadcaster.Middleware(func(fn broadcaster.HandlerFunc) broadcaster.HandlerFunc {
			return broadcaster.HandlerFunc(func(c []*broadcaster.Client, msgs []*broadcaster.Message) {
				for _, v := range msgs {
					(*v).(map[string]interface{})["author"] = fmt.Sprintf("written by %s", (*v).(map[string]interface{})["author"])
				}
				fn(c, msgs)
			})
		}),
	)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		panic("err: http.ListenAndServe failed for %s" + err.Error())
	}
}
