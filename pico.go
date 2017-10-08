package picoweb

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"strconv"

	"github.com/julienschmidt/httprouter"

	"strings"

	"github.com/googollee/go-socket.io"
)

var (
	RequestCount int
	isDev        bool
	flash        Flash
)

type Pico struct {
	mux          *httprouter.Router
	server       *http.Server
	c            chan os.Signal
	sio          *socketio.Server
	trackSession bool
	cookieName   string
	pre          PicoHandler
	post         PicoHandler
}

type PicoHandler func(c *Context)

func (p *Pico) Get(pattern string, fn PicoHandler) {
	p.mux.GET(pattern, middle(fn))
}

func (p *Pico) Post(pattern string, fn PicoHandler) {
	p.mux.POST(pattern, middle(fn))
}

func (p *Pico) Options(pattern string, fn PicoHandler) {
	p.mux.OPTIONS(pattern, middle(fn))
}

func (p *Pico) Put(pattern string, fn PicoHandler) {
	p.mux.PUT(pattern, middle(fn))
}

func (p *Pico) Delete(pattern string, fn PicoHandler) {
	p.mux.DELETE(pattern, middle(fn))
}

func (p *Pico) Static(urlPath, diskPath string) {
	if urlPath[len(urlPath)-1] == '/' {
		urlPath = urlPath[:len(urlPath)-1]
	}

	if urlPath[0:1] != "/" {
		urlPath = "/" + urlPath
	}

	p.mux.ServeFiles(urlPath+"/*filepath", http.Dir(diskPath))
}

func (p *Pico) EnableSocketIoOn(url string) {
	var err error
	p.sio, err = socketio.NewServer(nil)
	p.mux.GET(url, middlehttp(p.sio))
	fmt.Println(err)
}

func (p *Pico) OnConnection(fn func(s Socket)) {
	p.sio.On("connection", socket_middle(fn))
}

func (p *Pico) OnError(fn func(s Socket, e error)) {
	p.sio.On("error", socket_middle_error(fn))
}

func (p *Pico) On(event string, fn func(msg string)) {
	p.sio.On(event, fn)
}

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
		Handler: p.mux,
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
	signal.Notify(p.c, os.Interrupt)
	signal.Notify(p.c, os.Kill)

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
			fn()
		}

		os.Exit(0)
	}()
}

func (p *Pico) Stop() {
	fmt.Println("Shutting Down server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	err := p.server.Shutdown(ctx)
	if err != nil {
		fmt.Println(err)
	}
	cancel()

	//p.server.Stop(time.Second * 2)
	fmt.Println("Waiting on Stop Channel")
	flash.Clear()
	//<-p.server.StopChan()
}

func New() *Pico {
	isDev = true
	return &Pico{mux: httprouter.New()}
}
