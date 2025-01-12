package picoweb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func ID() string {
	u := uuid.New()
	return strings.Replace(u.String(), "-", "", -1)
}

type genericconnectionhandler struct {
	ID            string
	ex            *GoChannel
	c             *websocket.Conn
	out           *wsDataGoChannel
	isOpen        bool
	clienthandler WsHandler
}

func (wh *genericconnectionhandler) Terminate() {
	wh.ex.In(struct{}{})
}

func (wh *genericconnectionhandler) Dispose() {
	wh.c.Close()
	wh.ex.Close()
}

func (wh *genericconnectionhandler) handle(ctx context.Context, opendata WsData) {
	wh.isOpen = true

	defer func() {
		wh.ex.Close()
		wh.isOpen = false
	}()

	openResponse := wh.clienthandler(&WSArgs{ID: wh.ID, Channel: "ws", Command: "ws_open", Body: opendata})
	if openResponse != nil {
		if openResponse.Bool("close") {
			wh.clienthandler(&WSArgs{ID: wh.ID, Channel: "ws", Command: "ws_close", Body: WsData{"count": connections.count()}})
			return
		}
		wh.out.In(openResponse)
	}

	go func() {
		for wh.isOpen {
			_, body, err := wh.c.ReadMessage()
			if err != nil || len(body) == 0 {
				wh.ex.In(struct{}{})
				wh.isOpen = false
				return
			}

			var data WsData
			err = json.Unmarshal(body, &data)
			if err != nil {
				fmt.Println("Error in incoming", err)
				continue
			}

			response := wh.clienthandler(&WSArgs{
				Body:    data,
				ID:      wh.ID,
				Command: data.String("cmd"),
				Channel: data.String("channel"),
				Group:   data.String("group"),
				Account: data.String("account"),
			})

			if response != nil {
				if response.Bool("close") {
					wh.ex.In(struct{}{})
					wh.isOpen = false
					return
				}
				wh.out.In(response)
			}
		}
	}()

out:
	for wh.isOpen {
		select {
		case outgoing := <-wh.out.Out():
			err := wh.c.WriteMessage(websocket.TextMessage, []byte(outgoing.Json()))
			if err != nil {
				fmt.Println("Outgoing write", err)
				break out
			}
		case <-ctx.Done():
			// fmt.Println("Exiting CTX")
			break out
		case <-wh.ex.Out():
			// fmt.Println("Exiting EX")
			break out
		}
	}

	wh.isOpen = false
	wh.out.Close()
}

func NewGenericHandler(c *websocket.Conn) *genericconnectionhandler {
	return &genericconnectionhandler{
		ID:     ID(),
		c:      c,
		isOpen: true,
		ex:     Channel(5), //  make(chan struct{}, 2),
		out:    WsDataGoChannel(10),
	}
}
