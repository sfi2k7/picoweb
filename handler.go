package picoweb

import (
	"crypto/rand"
	"fmt"

	"github.com/gorilla/websocket"
)

type handler struct {
	ID          string
	c           *websocket.Conn
	ex          chan bool
	msgs        chan []byte
	isOpen      bool
	p           *Pico
	isConnected bool
}

// Note - NOT RFC4122 compliant
func UUID() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X%X%X%X%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}

func ID() string {
	return UUID()
	// u, _ := uuid.NewV4()
	// return strings.Replace(u.String(), "-", "", -1)
}

func (h *handler) init(c *websocket.Conn) string {
	h.ex = make(chan bool, 1)
	h.msgs = make(chan []byte, 10)
	h.c = c
	h.ID = ID()
	h.isOpen = true
	h.isConnected = true
	return h.ID
}

func (h *handler) forceExit() {
	h.ex <- true
}

func (h *handler) dispose() {
	//h.isOpen = false
	close(h.ex)
	close(h.msgs)
	//time.Sleep(time.Second * 10)
	h.c = nil
	if isDev == true {
		fmt.Println("Clean up Done")
	}
	h.isConnected = false
	h.isOpen = false
}

func (h *handler) handle() {
	defer func() {
		if data := recover(); data != nil {
			fmt.Println("Recover", data)
		}
	}()

	go func() {
		for {
			t, body, err := h.c.ReadMessage()

			if err != nil {
				h.ex <- true
				h.isConnected = false
				break
			}

			if isDev == true {
				fmt.Println("Got Message")
			}

			wshandler.Message(&WSMsgContext{
				MessageBody: body,
				MessageType: t,
				WSContext: &WSContext{
					conn:         h.c,
					ConnectionID: h.ID,
					p:            h.p,
				},
			})

			// if h.p.onMsg != nil {
			// 	go h.p.onMsg(&WSMsgContext{
			// 		MessageBody: body,
			// 		MessageType: t,
			// 		WSContext: &WSContext{
			// 			conn:         h.c,
			// 			ConnectionID: h.ID,
			// 			p:            h.p,
			// 		},
			// 	})
			// }
		}
	}()

	wshandler.Open(&WSContext{
		conn:         h.c,
		ConnectionID: h.ID,
		p:            h.p,
	})

	// if h.p.onConnect != nil {
	// 	h.p.onConnect(&WSContext{
	// 		conn:         h.c,
	// 		ConnectionID: h.ID,
	// 		p:            h.p,
	// 	})
	// }

	for {
		select {
		case msg := <-h.msgs:
			if msg == nil {
				return
			}

			err := h.c.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				h.ex <- true
			}
		case <-h.ex:
			h.isConnected = false
			return
		}
	}
}
