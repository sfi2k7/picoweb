# picoweb - tiny web franework

Teeny Tiny Web Wrapper around httprouter
# Features
- graceful shutdown
- socketio
- websocket
- more on the way

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