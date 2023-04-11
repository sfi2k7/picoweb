package main

import (
	"fmt"

	"github.com/sfi2k7/picoweb"
)

var (
	p *picoweb.Pico
)

func handler(args *picoweb.WSArgs) picoweb.WsData {
	fmt.Println("args", args)
	if args.Command == "hello" {
		return picoweb.WsData{"message": "Hello from server"}
	}

	return nil
}

// func background(ctx context.Context) {
// 	timer := time.NewTimer(5 * time.Second)

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case <-timer.C:
// 			// fmt.Println("background running")
// 			p.BroadcastWS(picoweb.WsData{"message": "Hello from server"})
// 			timer.Reset(5 * time.Second)
// 		}
// 	}
// }

func main() {
	p = picoweb.New()
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// go background(ctx)

	p.StopOnInt()
	p.CustomNotFound()
	p.SetAppName("picowebtest")

	p.WS("/ws", handler)
	p.Get("/", func(c *picoweb.Context) {
		c.View("./index.html", nil)
	})
	fmt.Println(p.Listen(57432))
}
