package picoweb

import "github.com/pkg/errors"

var (
	ErrChannelClosed = errors.New("channel closed")
)

type GoChannel struct {
	IsClosed bool
	c        chan interface{}
}

func (gc *GoChannel) In(v interface{}) error {
	if gc.IsClosed {
		// fmt.Println("Channel is closed")
		return ErrChannelClosed
	}

	gc.c <- v
	return nil
}

func (gc *GoChannel) Out() chan interface{} {
	return gc.c
}

func (gc *GoChannel) Close() {
	if gc.IsClosed {
		return
	}

	gc.IsClosed = true
	close(gc.c)
}

func Channel(cap int) *GoChannel {
	return &GoChannel{
		c: make(chan interface{}, cap),
	}
}

type wsDataGoChannel struct {
	IsClosed bool
	c        chan WsData
}

func (gc *wsDataGoChannel) In(v WsData) error {
	if gc.IsClosed {
		// fmt.Println("Channel is closed")
		return ErrChannelClosed
	}

	gc.c <- v
	return nil
}

func (gc *wsDataGoChannel) Out() chan WsData {
	return gc.c
}

func (gc *wsDataGoChannel) Close() {
	if gc.IsClosed {
		return
	}

	gc.IsClosed = true
	close(gc.c)
}

func WsDataGoChannel(cap int) *wsDataGoChannel {
	return &wsDataGoChannel{
		c: make(chan WsData, cap),
	}
}
