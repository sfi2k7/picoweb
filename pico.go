package picoweb

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"strconv"

	"github.com/julienschmidt/httprouter"
	mgo "gopkg.in/mgo.v2"

	"strings"
)

//var upgrader = websocket.Upgrader{EnableCompression: true, HandshakeTimeout: time.Second * 5, ReadBufferSize: 4096, WriteBufferSize: 4096}
var baseSession *mgo.Session

var (
	RequestCount  uint64
	isDev         bool
	flash         Flash
	mongoURL      string
	redisURL      string
	redisPassword string
)
var (
	startedOn         time.Time
	WSConnectionCount uint64
)

type Pico struct {
	Mux    *httprouter.Router
	server *http.Server
	c      chan os.Signal
	//sio          *socketio.Server
	trackSession bool
	cookieName   string
	pre          PicoHandler
	post         PicoHandler

	onMsg       func(c *WSMsgContext)
	onConnect   func(c *WSContext)
	onError     func(err error)
	onClose     func(c *WSContext)
	connections *mmap
}

type PicoHandler func(c *Context)

func (p *Pico) MongoURL(murl string) {
	mongoURL = murl
}

func (p *Pico) RedisURL(rurl string, redispassword ...string) {
	redisURL = rurl
	if len(redispassword) > 0 {
		redisPassword = redispassword[0]
	}
}

func (p *Pico) Get(pattern string, fn PicoHandler) {
	p.Mux.GET(pattern, middle(fn))
}

func (p *Pico) Post(pattern string, fn PicoHandler) {
	p.Mux.POST(pattern, middle(fn))
}

func (p *Pico) Options(pattern string, fn PicoHandler) {
	p.Mux.OPTIONS(pattern, middle(fn))
}

func (p *Pico) Put(pattern string, fn PicoHandler) {
	p.Mux.PUT(pattern, middle(fn))
}

func (p *Pico) Delete(pattern string, fn PicoHandler) {
	p.Mux.DELETE(pattern, middle(fn))
}

func (p *Pico) HandleWS(pattern string) {
	if p.connections != nil {
		return
	}

	p.connections = newmmap()
	p.Get(pattern, p.mainEndpoint)
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

// func (p *Pico) EnableSocketIoOn(url string) {
// 	var err error
// 	p.sio, err = socketio.NewServer(nil)
// 	p.Mux.GET(url, middlehttp(p.sio))
// 	fmt.Println(err)
// }

// func (p *Pico) OnConnection(fn func(s Socket)) {
// 	p.sio.On("connection", socket_middle(fn))
// }

// func (p *Pico) OnError(fn func(s Socket, e error)) {
// 	p.sio.On("error", socket_middle_error(fn))
// }

// func (p *Pico) On(event string, fn func(msg string)) {
// 	p.sio.On(event, fn)
// }

func (p *Pico) GetFlash(sessionId string) interface{} {
	return flash.Get(sessionId)
}

func (p *Pico) SetFlash(sessionId string, value interface{}) {
	flash.Set(sessionId, value)
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

func (p *Pico) ListenTLS(port, cert, key string) {
	cer, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		log.Println(err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	p.server = &http.Server{
		Addr:      port,
		Handler:   p.Mux,
		TLSConfig: config,
	}
	flash = make(Flash)
	p.server.ListenAndServeTLS(cert, key)
}

func (p *Pico) Listen(port int) error {
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

	return p.server.ListenAndServe()
}

func (p *Pico) Production() {
	isDev = false
}

func (p *Pico) StopOnInt() {
	p.StopOnIntWithFunc(nil)
}

func (p *Pico) StopOnIntWithFunc(fn func()) {
	p.c = make(chan os.Signal, 1)
	signal.Notify(p.c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	go func() {
		<-p.c

		if isDev {
			fmt.Println("Shutting Down!")
		}

		p.Stop()

		if isDev {
			fmt.Println("Done!")
		}

		close(p.c)

		if fn != nil {
			fmt.Println("Calling INT callback")
			fn()
		}

		fmt.Println("Exiting to OS")
		os.Exit(0)
	}()
}

func (p *Pico) Stop() {

	fmt.Println("Shutting Down server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := p.server.Shutdown(ctx)

	if err != nil {
		fmt.Println(err)
	}

	// cancel()

	//p.server.Stop(time.Second * 2)
	fmt.Println("Shutdown complete")
	flash.Clear()
	//<-p.server.StopChan()
}

func New() *Pico {
	isDev = true
	return &Pico{Mux: httprouter.New()}
}
