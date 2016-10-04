# go-websocket-json-broadcaster

This is based golang-samples/websocket

go-websocket-json-broadcaster is serve JSON from one client to every connected clients through on websocket.
It has hook middleware mechanism when broadcast JSON to clients like below. 

```go
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
```

It means message translation. 

```
[{"author":"hatajoe","body":"hello"}]
```

to

```
[{"author":"written by hatajoe","body":"hello"}]
```

Hook middleware can filter clients to broadcast depending on you. (but it's not sufficient yet)

# cmd

You can run go-websocket-json-broadcaster sample easily.

```
% cd cmd
% go run main.go -p 8080 -s 10 -e "/broadcast"
```

and connect to.

```
% wscat -c ws://127.0.0.1:8080/bloadcast -o localhost:8080
```

and write JSON array formatted string will broadcast other clients.

## options 

```
  -e string
    	websocket application endpoint path (default "/")
  -p int
    	websocket listen port (default 9218)
  -s int
    	message buffer size per client (default 100)
```

# LICENCE

Same as golang-samples/websocket


