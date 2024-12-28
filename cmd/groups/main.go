package main

import (
	"fmt"
	"net/http"

	"github.com/sfi2k7/picoweb"
)

func main() {
	/*
		//12/29/2024
		TODO: Able to define groups

		var router = picoweb.New()
		router.Use(commonmiddleware)

		api1 := router.Group("/api/v1")
		api1.Use(api1middleware)

		api2 := router.Group("/api/v2")
		api2.Use(api2middleware)

		api1.Get("/users", usersapiv1)
		api2.Get("/users", usersapiv2)

		router.Run(":8080")

	*/

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

	// //root middlewares
	// router.Use(func(c *picoweb.Context) bool {
	// 	fmt.Println("Middleware root")
	// 	return true
	// })

	// //root middlewares (must run)
	// router.Must(func(c *picoweb.Context) {
	// 	fmt.Println("Must middleware root")
	// })

	// //root Routes
	// router.Get("/", func(c *picoweb.Context) {
	// 	c.String("Hello World at /")
	// })

	// api := router.Group("/api")
	// api.GroupOptions().SkipMiddlewares()
	// //api middleware
	// api.Use(func(c *picoweb.Context) bool {
	// 	fmt.Println("Middleware api")
	// 	return true
	// })

	// api.Must(func(c *picoweb.Context) {
	// 	fmt.Println("Api must middleware (api logger)")
	// })

	// //Api routes
	// api.Get("/", func(c *picoweb.Context) {
	// 	c.String("Hello World at /api")
	// })

	// usersapi := api.Group("/users")
	// //users middleware
	// usersapi.Use(func(c *picoweb.Context) bool {
	// 	fmt.Println("Middleware users")
	// 	return true
	// })

	// //users routes
	// usersapi.Get("/", func(c *picoweb.Context) {
	// 	c.String("Hello World at /api/users")
	// })

	// // usersapi.GroupOptions().SkipMiddlewares()
	// _, _ = router.Ws("/ws", func(c *picoweb.WSArgs) picoweb.WsData {
	// 	fmt.Println("Websocket command", c.Command)
	// 	return picoweb.WsData{"response": "ok"}
	// })

}
