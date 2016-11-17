# picoweb

Teeny Tiny Web Wrapper arround GO's built in MUX
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
        pico.HandleFunc("/", Home)
        pico.Listen(7777)
    }

```