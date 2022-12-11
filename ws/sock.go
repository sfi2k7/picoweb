package ws

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/sfi2k7/picoweb"
)

type WsHandlers interface {
	OnOpen(c *WSContext)
	OnClose(c *WSContext)
	OnData(c *WSMsgContext)
}

// type WsHandler interface {
// 	Open(c *WSContext)
// 	Message(c *WSMsgContext)
// 	Close(c *WSContext)
// 	Error(error)
// }

var (
	isWsSet           bool
	WSConnectionCount uint64
)

type WSConn websocket.Conn

func MainEndpoint(c *picoweb.Context) {
	atomic.AddUint64(&WSConnectionCount, 1)
	con, err := c.Upgrade()

	if err != nil {
		c.Status(http.StatusBadRequest)
		//TODO: ERROR handler
		return
	}

	fmt.Println("WS Count", atomic.LoadUint64(&WSConnectionCount))

	//TODO: Open Handler

	wsLoop(c, con)

	atomic.AddUint64(&WSConnectionCount, ^uint64(0))
}

func wsLoop(c *picoweb.Context, con *websocket.Conn) {
	h := &handler{ctx: c}
	id := h.init(con)

	connections.add(id, h)

	defer func() {
		connections.remove(id)

		//TODO: Close Handler
		// wshandler.Close()
		handlers.OnClose(&WSContext{conn: con, ConnectionID: id})
		h.dispose()
		h.isConnected = false
	}()

	h.handle()
}

// func (p *Pico) SendWS(cid string, data interface{}) error {
// 	h := p.connections.get(cid)
// 	if h == nil {
// 		return errors.New("Connection not found")
// 	}

// 	h.msgs <- data
// 	return nil
// }

// func (p *Pico) SendJson(cid string, o interface{}) error {
// 	jsoned, _ := json.Marshal(o)
// 	return p.SendWS(cid, jsoned)
// }
