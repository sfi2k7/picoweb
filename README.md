## WIP - Not ready to use
# picoweb - tiny web franework

Teeny Tiny Web Wrapper around httprouter

## Features
- Fast - thanks to _httprouter_
- **graceful** shutdown
- socketio
- raw WebSocket (soon)
- more on the way

```GO
    package main

    import (
        "fmt"

        "github.com/sfi2k7/picoweb"
    )

    func Home(c *picoweb.Context) {
        fmt.Fprint(c, "Hello Pico")
    }

    func main() {
        pico := picoweb.New()
        pico.Static("/static", "./static")
        pico.Get("/", Home)
        pico.Listen(7777)
    }
```

### SocketIO - Removes in latest Version

```GO

    pico.EnableSocketIoOn("/socket.io/")
    pico.OnConnection(func(s socketio.Socket) {
        s.emit("welcome","Welcome to Pico web framewwork")
    })

```

### Enable Production Mode

```GO
    pico.Production()
```

### Parameters

```GO
    pico.Get("/:name",func (c *picoweb.Context){
        userName := c.Params("name")
    })
```

### License
MIT - Please see the `LICENSE` file
