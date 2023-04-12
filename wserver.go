package picoweb

import (
	"context"
	"fmt"
	"net/http"
)

var (
	connections *genericmmap
)

type WsHandler func(args *WSArgs) WsData

type WSArgs struct {
	ID      string
	Command string
	Body    WsData
	Channel string
	Group   string
	Node    string
	Account string
}

type GenericWsGoServer struct {
	isShuttingDown bool
	MessageHandler WsHandler
}

func (wsg *GenericWsGoServer) Close() {
	wsg.isShuttingDown = true
	connections.closeAll()
}

func (wsg *GenericWsGoServer) Handle(c *Context) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovering from", r)
		}
	}()

	if wsg.isShuttingDown {
		c.WriteHeader(http.StatusInternalServerError)
		return
	}

	con, err := c.Upgrade()
	if err != nil {
		c.WriteHeader(http.StatusBadRequest)
		return
	}

	var openData = WsData{}

	for k, q := range c.URL().Query() {
		if len(q) == 0 {
			continue
		}
		openData[k] = q[0]
	}

	for k, v := range c.r.Header {
		if len(v) == 0 {
			continue
		}

		openData[k] = v[0]
	}

	handler := NewGenericHandler(con)
	handler.clienthandler = wsg.MessageHandler
	defer handler.Dispose()

	connections.add(handler.ID, handler)

	wsg.MessageHandler(&WSArgs{ID: handler.ID, Channel: "ws", Command: "ws_open", Body: openData})
	wsg.MessageHandler(&WSArgs{ID: handler.ID, Channel: "ws", Command: "ws_count", Body: WsData{"count": connections.count()}})

	handler.handle(context.Background())

	connections.remove(handler.ID)

	wsg.MessageHandler(&WSArgs{ID: handler.ID, Channel: "ws", Command: "ws_close"})
	wsg.MessageHandler(&WSArgs{ID: handler.ID, Channel: "ws", Command: "ws_remove"})
	wsg.MessageHandler(&WSArgs{ID: handler.ID, Channel: "ws", Command: "ws_count", Body: WsData{"count": connections.count()}})

	handler.clienthandler = nil
}
