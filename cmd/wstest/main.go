package main

import (
	"fmt"
	"net/url"
	"runtime"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sfi2k7/picoweb"
)

var (
	web *picoweb.Pico
	wg  sync.WaitGroup
)

func onOpen(c *picoweb.WSContext) {
	c.Text("Hello " + c.ConnectionID)
}

func onClose(c *picoweb.WSContext) {

}

func onMsg(c *picoweb.WSMsgContext) {

}

func ws_server() {
	web = picoweb.New()
	web.HandleWS("/ws")
	web.OnWSOpen(onOpen)
	web.OnWSClose(onClose)
	web.OnWSMsg(onMsg)
	web.Production()
	web.Listen(9562)
}

func ws_client() {
	defer wg.Done()

	u := url.URL{Scheme: "ws", Host: "localhost:9562", Path: "/ws"}
	//log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	_, _, err = c.ReadMessage()
	if err != nil {
		return
	}

	//fmt.Println("From Server", string(msg))
	//time.Sleep(time.Second * 1)
}

func main() {
	go ws_server()

	defer web.Stop()
	wg = sync.WaitGroup{}
	y := 0
	var memUsage runtime.MemStats
	for {
		runtime.ReadMemStats(&memUsage)
		fmt.Println("Batch", strconv.Itoa(y), "GO", runtime.NumGoroutine(), "Alloc", memUsage.Alloc/(1024*1024), "Live", memUsage.Mallocs-memUsage.Frees, "Pause", memUsage.PauseTotalNs)
		for x := 0; x < 200; x++ {
			//fmt.Println("Starting", strconv.Itoa(x))
			go ws_client()
			wg.Add(1)
		}
		wg.Wait()

		y++

		if y > 100 {
			break
		}
	}

}
