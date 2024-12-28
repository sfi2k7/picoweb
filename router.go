package picoweb

import (
	"fmt"
	"path"
)

type Router struct {
	mux             *Pico
	prefix          string
	parent          *Router
	middlewares     []PicoMiddleWareHandler
	mustmiddlewares []PicoHandler
	port            int
	cert            string
	key             string
	so              *serveroptions
	gopt            *serveroptions
}

// TODO Allow same options per Group
type serveroptions struct {
	skipmiddlewares bool
	skipmusts       bool
}

type Config struct {
	r *Router
}

type GroupOptions struct {
	r *Router
}

// Config gets the config for the server
func (r *Router) Config() *Config {
	if r.parent != nil {
		panic("Config can only be called on root router")
	}

	return &Config{r: r}
}

// BroadcastWs sends a message to all websockets
// data is the data to send
// exclude is a list of ids to exclude
func (r *Router) BroadcaseWs(data WsData, exclude ...string) {
	r.mux.BroadcastWS(data, exclude...)
}

// SendWs sends a message to a websocket
// id is the id of the websocket
// data is the data to send
func (r *Router) SendWs(id string, data WsData) {
	r.mux.SendWS(id, data)
}

// GroupOptions allows you to set options for a group of routes
func (r *Router) GroupOptions() *GroupOptions {
	if r.parent == nil {
		panic("Group Options can only be called on a group router")
	}

	return &GroupOptions{r: r}
}

// SkipMiddlewares skips all middlewares
// Middlewares are functions that run before the route handler
func (g *GroupOptions) SkipMiddlewares() *GroupOptions {
	g.r.gopt.skipmiddlewares = true
	return g
}

// SkipMusts skips all must middlewares
// Must middlewares are middlewares that must run after the route handler
func (g *GroupOptions) SkipMusts() *GroupOptions {
	g.r.gopt.skipmusts = true
	return g
}

// SetDev sets the server to development mode
func (c *Config) SetDev(dev bool) *Config {
	if !dev {
		c.r.mux.Production()
	}
	return c
}

// SetPort sets the port for the server
func (c *Config) SetPort(port int) *Config {
	c.r.port = port
	return c
}

// StopOnIntrupt stops the server on interrupt signal
func (c *Config) StopOnInterrupt() *Config {
	c.r.mux.StopOnInt()
	return c
}

// StopOnIntruptWithFunc stops the server on interrupt signal and runs a function
// fn is the function to run
func (c *Config) StopOnInterruptWithFunc(fn func()) *Config {
	c.r.mux.StopOnIntWithFunc(fn)
	return c
}

// SkipAllMiddlewares skips all middlewares
// Middlewares are functions that run before the route handler
func (c *Config) SkipAllMiddlewares() *Config {
	c.r.so.skipmiddlewares = true
	return c
}

// Static sets a static file server
// urlPath is the path to serve the files
// diskPath is the path to the files on disk
func (c *Config) Static(urlPath, diskPath string) *Config {
	c.r.mux.Static(urlPath, diskPath)
	return c
}

// UseSSL sets the server to use SSL
// cert and key are the paths to the certificate and key files
func (c *Config) UseSSL(cert, key string) *Config {
	c.r.cert = cert
	c.r.key = key
	return c
}

// SkipMusts skips all must middlewares
// Must middlewares are middlewares that must run after the route handler
func (c *Config) SkipMusts() *Config {
	c.r.so.skipmusts = true
	return c
}

// GlobalOPTIONS sets the handler for global OPTIONS requests
// This is the same as setting a route for the path with the method OPTIONS
func (c *Config) GlobalOPTIONS(fn PicoHandler) *Config {
	c.r.mux.Mux.GlobalOPTIONS = picohandlertohttphandler(fn)
	return c
}

// HandleOPTIONS sets the server to handle OPTIONS requests
func (c *Config) HandleOPTIONS() *Config {
	c.r.mux.Mux.HandleOPTIONS = true
	return c
}

// MethodNotAllowed sets the handler for when a method is not allowed
// This is the same as setting a route for the path with the method not allowed
func (c *Config) MethodNotAllowed(fn PicoHandler) *Config {
	c.r.mux.Mux.MethodNotAllowed = picohandlertohttphandler(fn)
	return c
}

// NotFound sets the handler for when a route is not found
// This is the same as setting a route for the path not found
func (c *Config) NotFound(fn PicoHandler) *Config {
	c.r.mux.Mux.NotFound = picohandlertohttphandler(fn)
	return c
}

// RedirectFixedPath sets the server to redirect fixed paths
// This is the same as setting a route for the path with the fixed path
func (c *Config) RedirectFixedPath() *Config {
	c.r.mux.Mux.RedirectFixedPath = true
	return c
}

// RedirectTrailingSlash sets the server to redirect trailing slashes
// This is the same as setting a route for the path with the trailing slash
func (c *Config) RedirectTrailingSlash() *Config {
	c.r.mux.Mux.RedirectTrailingSlash = true
	return c
}

// Group sets the prefix for a group of routes
// prefix is the prefix for the group
// returns a new router
func (r *Router) Group(prefix string) *Router {
	router := &Router{so: r.so, gopt: &serveroptions{}, parent: r, mux: r.mux, prefix: path.Join(r.prefix, prefix)}
	return router
}

// Ws sets a websocket endpoint
// pattern is the path for the websocket
// fn is the handler for the websocket
// returns a broadcast function and a send function
func (r *Router) Ws(pattern string, fn WsHandler) (broadcase func(data WsData, exclude ...string), send func(id string, data WsData)) {
	if r.parent != nil {
		panic("Websocket endpoint can only be defined at root level")
	}

	r.mux.Ws(pattern, fn)
	return r.mux.BroadcastWS, r.mux.SendWS
}

func (r *Router) WsSimple(pattern string, fn WsHandler) {
	if r.parent != nil {
		panic("Websocket endpoint can only be defined at root level")
	}

	r.mux.Ws(pattern, fn)
}

// Get sets a GET route
// pattern is the path for the route
// fn is the handler for the route
func (r *Router) Get(pattern string, fn PicoHandler) {
	r.mux.Get(path.Join(r.prefix, pattern), r.middleware(fn))
}

// Post sets a POST route
// pattern is the path for the route
// fn is the handler for the route
func (r *Router) Post(pattern string, fn PicoHandler) {
	r.mux.Post(path.Join(r.prefix, pattern), r.middleware(fn))
}

// Put sets a PUT route
// pattern is the path for the route
func (r *Router) Put(pattern string, fn PicoHandler) {
	r.mux.Put(path.Join(r.prefix, pattern), r.middleware(fn))
}

// Delete sets a DELETE route
// pattern is the path for the route
func (r *Router) Delete(pattern string, fn PicoHandler) {
	r.mux.Delete(path.Join(r.prefix, pattern), r.middleware(fn))
}

// Patch sets a PATCH route
// pattern is the path for the route
func (r *Router) Patch(pattern string, fn PicoHandler) {
	r.mux.Patch(path.Join(r.prefix, pattern), r.middleware(fn))
}

// Options sets a OPTIONS route
// pattern is the path for the route
func (r *Router) Options(pattern string, fn PicoHandler) {
	r.mux.Options(path.Join(r.prefix, pattern), r.middleware(fn))
}

// Must sets a must middleware that must run after the route handler
// fn is the must middleware
func (r *Router) Must(fn PicoHandler) {
	r.mustmiddlewares = append(r.mustmiddlewares, fn)
}

// Use sets a middleware
// fn is the middleware
func (r *Router) Use(fn PicoMiddleWareHandler) {
	r.middlewares = append(r.middlewares, fn)
}

func (r *Router) runMust(c *Context) {
	//Bottom Up - run parent must middlewares last
	if !r.gopt.skipmusts {
		for _, middle := range r.mustmiddlewares {
			middle(c)
		}
	}

	if r.parent != nil && !r.so.skipmusts {
		r.parent.runMust(c)
	}
}

func (r *Router) runMiddlewares(c *Context) bool {
	// fmt.Println("r.so", r.so, "r.gopt", r.gopt)

	//TOP DOWN - run parent middlewares first
	if r.parent != nil && !r.so.skipmiddlewares {
		if !r.parent.runMiddlewares(c) {
			return false
		}
	}

	if r.gopt.skipmiddlewares {
		return true
	}

	for _, middle := range r.middlewares {
		if !middle(c) {
			return false
		}
	}

	return true
}

func (r *Router) middleware(fn PicoHandler) PicoHandler {
	return func(c *Context) {
		movenext := r.runMiddlewares(c)

		if movenext || r.gopt.skipmiddlewares {
			fn(c)
		}

		r.runMust(c)
	}
}

// NewRouter creates a new router
// returns a new router
func NewRouter() *Router {
	return &Router{gopt: &serveroptions{}, so: &serveroptions{}, parent: nil, port: 8080, mux: New()}
}

// StartServer starts the server
// returns an error if the server fails to start
func (r *Router) StartServer() error {
	if len(r.cert) > 0 && len(r.key) > 0 {
		return r.mux.ListenTLS(":"+fmt.Sprint(r.port), r.cert, r.key)
	}

	return r.mux.Listen(r.port)
}

// StopServer stops the server
// returns an error if the server fails to stop
func (r *Router) StopServer() error {
	return r.mux.Stop()
}
