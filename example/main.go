package main

import (
	"fmt"

	"github.com/sfi2k7/picoweb"
)

func main() {
	p := picoweb.New()
	p.StopOnInt()
	p.Get("/", func(c *picoweb.Context) {
		c.View("./index.html", struct{ Name string }{Name: "Faisal"})
		//c.String("Hello World")
	})
	fmt.Println(p.Listen(57432))
}
