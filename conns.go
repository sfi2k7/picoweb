package picoweb

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	locker = sync.Mutex{}
)

/*
	onOpen:
		connections.add(cid)
	onLogin:
		u2cconnection[cid] = uid
		uconnection[uid].add(cid)

	onClose:
		connections.remove(cid)
		uname = u2cconnections[cid]
		u2cconnections.remove(cid)
		uconnection[uname].remove(cid)

*/

type uconn map[string][]*websocket.Conn
type conns map[*websocket.Conn]string

func add(c *websocket.Conn, u string) {

}
