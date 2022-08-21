package picoweb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type WsHandler interface {
	Open(c *WSContext)
	Message(c *WSMsgContext)
	Close(c *WSContext)
	Error(error)
}

var (
	isWsSet bool
)

type WSConn websocket.Conn

func (p *Pico) OnWSMsg(fn func(c *WSMsgContext)) {
	p.onMsg = fn
}

func (p *Pico) OnWSOpen(fn func(c *WSContext)) {
	p.onConnect = fn
}

func (p *Pico) OnWSClose(fn func(c *WSContext)) {
	p.onClose = fn
}

func (p *Pico) OnWSError(fn func(err error)) {
	p.onError = fn
}

func (p *Pico) CloseWSConn(cid string) {
	h := p.connections.get(cid)
	if h == nil {
		return
	}
	h.forceExit()
}

func (p *Pico) mainEndpoint(c *Context) {
	// var memUsage runtime.MemStats
	// if isDev{
	// 	runtime.ReadMemStats(&memUsage)
	// }

	atomic.AddUint64(&WSConnectionCount, 1)

	// if isDev{
	// 	fmt.Println("Alloc", memUsage.Alloc/1024*1024, "Live", memUsage.Mallocs-memUsage.Frees)
	// }

	if isDev {
		fmt.Println("Go Routine Count", runtime.NumGoroutine())
	}

	con, err := c.Upgrade()

	if err != nil {
		c.Status(http.StatusBadRequest)
		wshandler.Error(err)
		// if p.onError != nil {
		// 	p.onError(err)
		// }
		return
	}

	if isDev {
		fmt.Println("WS Count", atomic.LoadUint64(&WSConnectionCount))
	}

	p.wsLoop(con)

	atomic.AddUint64(&WSConnectionCount, ^uint64(0))
}

func (p *Pico) wsLoop(con *websocket.Conn) {
	h := &handler{p: p}
	id := h.init(con)

	if isDev {
		fmt.Println("Initing WS conn", id)
	}

	p.connections.add(id, h)

	defer func() {
		p.connections.remove(id)

		// if p.onClose != nil {
		// 	p.onClose()
		// }
		wshandler.Close(&WSContext{p: p, conn: con, ConnectionID: id})

		h.dispose()
		h.isConnected = false
		if isDev {
			fmt.Println("Connection Closed and Disposed", p.connections.count())
		}
		//h = nil
	}()

	if isDev {
		fmt.Println("Connection Made", p.connections.count())
	}

	h.handle()
}

func (p *Pico) SendWS(cid string, data interface{}) error {
	h := p.connections.get(cid)
	if h == nil {
		return errors.New("Connection not found")
	}

	h.msgs <- data
	return nil
}

func (p *Pico) SendJson(cid string, o interface{}) error {
	jsoned, _ := json.Marshal(o)
	return p.SendWS(cid, jsoned)
}
