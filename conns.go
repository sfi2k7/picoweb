package picoweb

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	locker = sync.Mutex{}
)

type uconn map[string][]*websocket.Conn
type conns map[*websocket.Conn]string

func add(c *websocket.Conn, u string) {

}
