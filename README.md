# picoweb

Teeny Tiny Web Wrapper arround httprouter
Uses github.com/tylerb/graceful as Server


```GO
   package main

    import (
        "fmt"

        "github.com/sfi2k7/picoweb"
    )

    func Home(c *picoweb.Context) {
        fmt.Fprint(c, "Hello World")
    }

    func main() {
        pico := picoweb.New()
        pico.Static("/static", "./static")
        pico.Get("/", Home)
        pico.Listen(7777)
    }


```