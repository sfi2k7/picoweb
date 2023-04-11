package main

import (
	"fmt"

	"github.com/sfi2k7/picoweb"
)

func handler(args *picoweb.WSArgs) picoweb.WsData {
	fmt.Println("args", args)
	if args.Command == "hello" {
		return picoweb.WsData{"message": "Hello from server"}
	}

	return nil
}

func main() {
	p := picoweb.New()

	p.StopOnInt()
	p.CustomNotFound()
	p.SetAppName("picowebtest")

	p.WS("/ws", handler)
	p.Get("/", func(c *picoweb.Context) {
		c.View("./index.html", nil)
	})
	fmt.Println(p.Listen(57432))
}
