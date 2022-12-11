package ws

import (
	"github.com/gorilla/websocket"
	"github.com/sfi2k7/picoweb"
)

type WSContext struct {
	*picoweb.Context
	ConnectionID string
	conn         *websocket.Conn
}

func (wsc *WSContext) WSSendText(msg string) error {
	return wsc.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (wsc *WSContext) WSSendJson(o interface{}) error {
	return wsc.conn.WriteJSON(o)
}

type WSMsgContext struct {
	*WSContext
	MessageType int
	MessageBody []byte
}
