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

### Latest Examples:

```GO
    router := picoweb.NewRouter()

	//WS Router
	router.WsSimple("/ws", func(args *picoweb.WSArgs) picoweb.WsData {
		if args.Command == "ws_open" {
			router.BroadcaseWs(picoweb.WsData{"response": "use joined"}, args.ID)
			return picoweb.WsData{"message": "welcome!"}
		}

		if args.Command == "ws_close" {
			router.BroadcaseWs(picoweb.WsData{"response": "user left"}, args.ID)
			fmt.Println("Users Count:", args.Body.Get("count"))
			return picoweb.WsData{"message": "goodbye!"}
		}

		router.SendWs(args.ID, picoweb.WsData{"response": "id - private"})
		return nil
	})

	//Root Routes
	router.Use(func(c *picoweb.Context) bool {
		fmt.Println("Middleware root")
		return true
	})

	//root middlewares (must run)
	router.Must(func(c *picoweb.Context) {
		fmt.Println("Must middleware root")
	})

	router.Get("/", func(c *picoweb.Context) {
		c.String("Hello World at /")
	})

	//API Routes
	api := router.Group("/api")
	api.GroupOptions().SkipMiddlewares()

	api.Use(func(c *picoweb.Context) bool {
		fmt.Println("Middleware api")
		return true
	})

	api.Must(func(c *picoweb.Context) {
		fmt.Println("Api must middleware (api logger)")
	})

	api.Get("/", func(c *picoweb.Context) {
		c.String("Hello World at /api")
	})

	//Api/Users Routes
	usersapi := api.Group("/users")

	usersapi.Use(func(c *picoweb.Context) bool {
		fmt.Println("Middleware /api/users")
		return true
	})

	usersapi.Must(func(c *picoweb.Context) {
		fmt.Println("Must middleware /api/users")
	})

	usersapi.Get("/", func(c *picoweb.Context) {
		c.View("./index.html", nil)
		// c.String("Hello World at /api/users")
	})

	config := router.Config()
	config.SetDev(true).StopOnInterrupt()
	config.Static("/static/", "./public")
	config.NotFound(func(c *picoweb.Context) {
		c.StatusWithString(http.StatusNotFound, "Route not defined:"+c.URL().Path)
	})
	config.SetPort(7865)

	router.StartServer()
```



### License
MIT - Please see the `LICENSE` file
