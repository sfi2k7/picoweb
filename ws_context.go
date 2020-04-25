package picoweb

import "github.com/gorilla/websocket"

type WSContext struct {
	ConnectionID string
	conn         *websocket.Conn
	p            *Pico
}

func (wsc *WSContext) Text(msg string) error {
	return wsc.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (wsc *WSContext) Json(o interface{}) error {
	return wsc.conn.WriteJSON(o)
}

func (wsc *WSContext) Login(u string) {
	wsc.p.connections.login(u, wsc.ConnectionID)
}

type WSMsgContext struct {
	*WSContext
	MessageType int
	MessageBody []byte
}
