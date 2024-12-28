package picoweb

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	mgo "gopkg.in/mgo.v2"

	"strings"
)

// var upgrader = websocket.Upgrader{EnableCompression: true, HandshakeTimeout: time.Second * 5, ReadBufferSize: 4096, WriteBufferSize: 4096}
var baseSession *mgo.Session
var skipmiddlewares bool
var skipmusts bool

var server *GenericWsGoServer

var (
	requestCount uint64
	isDev        bool
	flash        Flash
	// mongoURL      string
	// redisURL      string
	// redisPassword string
)

var (
	startedOn time.Time
)

type Pico struct {
	Mux    *httprouter.Router
	server *http.Server
	c      chan os.Signal
	//sio          *socketio.Server
	// trackSession bool
	// pre          PicoHandler
	// post         PicoHandler
	useAppManager bool
	appName       string
}

type PicoHandler func(c *Context)

// type PicoHandler func(c *Context) interface{}

// func (p *Pico) MongoURL(murl string) {
// 	mongoURL = murl
// }

func (p *Pico) SendWS(id string, data WsData) {
	if server == nil {
		return
	}

	connections.send(id, data)
}

func (p *Pico) BroadcastWS(data WsData, exclude ...string) {
	if server == nil {
		return
	}

	connections.broadcast(data, exclude...)
}

// func (p *Pico) RedisURL(rurl string, redispassword ...string) {
// 	redisURL = rurl
// 	if len(redispassword) > 0 {
// 		redisPassword = redispassword[0]
// 	}
// }

func (p *Pico) Get(pattern string, fn PicoHandler) {
	p.Mux.GET(pattern, middle(fn, p.appName))
}

func (p *Pico) Ws(pattern string, mh WsHandler) {
	if server != nil {
		panic(errors.New("only one websocket server is allowed per application"))
	}

	if mh == nil {
		panic(errors.New("websocket handler cannot be nil"))
	}

	server = &GenericWsGoServer{MessageHandler: mh}
	connections = newgenericmmap()
	p.Mux.GET(pattern, middle(server.Handle, p.appName))
}

func (p *Pico) Post(pattern string, fn PicoHandler) {
	p.Mux.POST(pattern, middle(fn, p.appName))
}

func (p *Pico) Options(pattern string, fn PicoHandler) {
	p.Mux.OPTIONS(pattern, middle(fn, p.appName))
}

func (p *Pico) Put(pattern string, fn PicoHandler) {
	p.Mux.PUT(pattern, middle(fn, p.appName))
}

func (p *Pico) Patch(pattern string, fn PicoHandler) {
	p.Mux.PATCH(pattern, middle(fn, p.appName))
}

func (p *Pico) Delete(pattern string, fn PicoHandler) {
	p.Mux.DELETE(pattern, middle(fn, p.appName))
}

func (p *Pico) StaticDefault(diskPath string) {
	p.Mux.ServeFiles("/*filepath", http.Dir(diskPath))
}

func (p *Pico) Static(urlPath, diskPath string) {
	if urlPath[len(urlPath)-1] == '/' {
		urlPath = urlPath[:len(urlPath)-1]
	}

	if urlPath[0:1] != "/" {
		urlPath = "/" + urlPath
	}

	p.Mux.ServeFiles(urlPath+"/*filepath", http.Dir(diskPath))
}

func (p *Pico) SkipAllMiddlewares() {
	skipmiddlewares = true
}

func (p *Pico) Before(m PicoMiddleWareHandler) {
	premiddlewares = append(premiddlewares, m)
}

func (p *Pico) Use(m PicoMiddleWareHandler) {
	middlewares = append(middlewares, m)
}

func (p *Pico) Must(m PicoMiddleWareHandler) {
	must = append(must, m)
}

func (p *Pico) After(m PicoMiddleWareHandler) {
	postmiddlewares = append(postmiddlewares, m)
}

func (p *Pico) GetFlash(sessionId string) interface{} {
	return flash.Get(sessionId)
}

func (p *Pico) SetFlash(sessionId string, value interface{}) {
	flash.Set(sessionId, value)
}

// obsolete
func (p *Pico) UseUserManager(url, password string) {
	// p.useAppManager = true
	// useUserManager(url, password)
}

// func (p *Pico) TrackSessionUsingCookie(cookieName string) {
// 	p.trackSession = true
// 	p.cookieName = cookieName
// }

func (p *Pico) ListenS(port string) {
	if strings.Index(port, ":") == 0 {
		port = port[1:]
	}

	po, err := strconv.ParseInt(port, 10, 32)
	if err != nil {
		panic("Port Error (:9999 etc)")
	}
	if po < 0 || po > 65000 {
		panic("PORT is OUT of Range")
	}
	p.Listen(int(po))
}

func (p *Pico) ListenTLS(port, cert, key string) error {
	cer, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		log.Println(err)
		return err
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	p.server = &http.Server{
		Addr:      port,
		Handler:   p.Mux,
		TLSConfig: config,
	}
	flash = make(Flash)
	return p.server.ListenAndServeTLS(cert, key)
}

func (p *Pico) Listen(port int) error {
	defer func() {
		if connections != nil {
			connections.closeAll()
		}
	}()

	envPort := os.Getenv("PORT")
	if len(envPort) > 0 {
		pi, err := strconv.Atoi(envPort)
		if err == nil && pi < 65000 {
			port = pi
		}
	}

	p.server = &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: p.Mux,
	}

	// p.server = &graceful.Server{
	// 	Timeout: 2 * time.Second,
	// 	Server:  &http.Server{},
	// }
	flash = make(Flash)
	if isDev {
		fmt.Println("Listing on " + strconv.Itoa(port))
	}

	rqc = reqcount{s: sync.Mutex{}}
	return p.server.ListenAndServe()
}

func (p *Pico) Production() {
	isDev = false
}

func (p *Pico) StopOnInt() {
	p.StopOnIntWithFunc(nil)
}

func (p *Pico) SetAppName(appname string) {
	p.appName = appname
}

func (p *Pico) CustomNotFound() {
	p.Mux.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(string("404 page not found " + p.appName)))
	})
}

func (p *Pico) StopOnIntWithFunc(fn func()) {
	p.c = make(chan os.Signal, 1)
	signal.Notify(p.c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-p.c

		if server != nil {
			server.isShuttingDown = true
			connections.closeAll()
		}

		if isDev {
			fmt.Println("Shutting Down!")
		}

		p.Stop()

		if isDev {
			fmt.Println("Done!")
		}

		if fn != nil {
			fmt.Println("Calling INT callback")
			fn()
		}

		close(p.c)

		fmt.Println("Exiting to OS")
		os.Exit(0)
	}()
}

func (p *Pico) Stop() error {

	fmt.Println("Shutting Down server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := p.server.Shutdown(ctx)

	flash.Clear()
	if err != nil {
		return err
	}
	fmt.Println("Shutdown complete")
	return nil
}

func New() *Pico {
	isDev = true
	return &Pico{Mux: httprouter.New()}
}

func picohandlertohttphandler(c PicoHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c(&Context{w: w, r: r})
	})
}
