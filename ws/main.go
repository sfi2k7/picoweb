package ws

import (
	"errors"

	"github.com/gorilla/websocket"
	"github.com/sfi2k7/picoweb"
)

var connections *mmap
var handlers WsHandlers

func init() {
	connections = newmmap()
}

func SetHandlers(hs WsHandlers) {
	handlers = hs
}

func SendTextToConnection(id string, data []byte) error {
	h := connections.get(id)
	if h == nil || !h.isConnected || !h.isOpen {
		return errors.New("client already disconnected")
	}
	return h.c.WriteMessage(websocket.TextMessage, data)
}

func SendJSONToConnection(id string, data interface{}) error {
	h := connections.get(id)
	if h == nil || !h.isConnected || !h.isOpen {
		return errors.New("client already disconnected")
	}
	return h.c.WriteJSON(data)
}

func WsHandler(c *picoweb.Context) {
	c.IsWebsocket = true
	MainEndpoint(c)
}
