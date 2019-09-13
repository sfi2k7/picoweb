package picoweb

import (
	"fmt"
	"net/http"
	"runtime"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var ()

var (
	isWsSet bool
)

type WSConn websocket.Conn
type WSContext struct {
	ConnectionID string
	Conn         *websocket.Conn
}

type WSMsgContext struct {
	*WSContext
	MessageType int
	MessageBody []byte
}

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
	atomic.AddUint64(&WSConnectionCount, 1)

	if isDev == true {
		fmt.Println("Go Count", runtime.NumGoroutine())
	}

	con, err := c.Upgrade()

	if err != nil {
		c.Status(http.StatusBadRequest)
		if p.onError != nil {
			p.onError(err)
		}
		return
	}

	if isDev == true {
		fmt.Println("WS Count", atomic.LoadUint64(&WSConnectionCount))
	}

	p.wsLoop(con)
}

func (p *Pico) wsLoop(con *websocket.Conn) {
	h := &handler{p: p}
	id := h.init(con)

	if isDev == true {
		fmt.Println("Initing WS conn", id)
	}

	p.connections.add(id, h)

	defer func() {
		p.connections.remove(id)

		if p.onClose != nil {
			p.onClose(&WSContext{Conn: con, ConnectionID: id})
		}

		if isDev == true {
			fmt.Println("Connection Closed", p.connections.count())
		}

		h.dispose()
	}()

	if isDev == true {
		fmt.Println("Connection Made", p.connections.count())
	}

	h.handle()
}

func (p *Pico) SendWS(cid string, body []byte) error {
	h := p.connections.get(cid)
	if h == nil {
		return errors.New("Connection not found")
	}
	h.msgs <- body
	return nil
}
