package picoweb

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tylerb/graceful"

	"github.com/googollee/go-socket.io"
)

var (
	RequestCount int
	isDev bool
)

type Pico struct {
	mux    *httprouter.Router
	server *graceful.Server
	c      chan os.Signal
	sio    *socketio.Server
}

type PicoHandler func(c *Context)

func (p *Pico) Get(pattern string, fn PicoHandler) {
	p.mux.GET(pattern, middle(fn))
}

func (p *Pico) Post(pattern string, fn PicoHandler) {
	p.mux.POST(pattern, middle(fn))
}

func (p *Pico) Put(pattern string, fn PicoHandler) {
	p.mux.PUT(pattern, middle(fn))
}

func (p *Pico) Delete(pattern string, fn PicoHandler) {
	p.mux.DELETE(pattern, middle(fn))
}

func (p *Pico) Static(urlPath, diskPath string) {
	p.mux.ServeFiles(urlPath+"/*filepath", http.Dir(diskPath))
}

func (p *Pico) EnableSocketIoOn(url string) {
	var err error
	p.sio, err = socketio.NewServer(nil)
	p.mux.GET(url, middlehttp(p.sio))
	fmt.Println(err)
}

func (p *Pico) OnConnection(fn func(s socketio.Socket)) {
	p.sio.On("connection", fn)
}

func (p *Pico) OnError(fn func(s socketio.Socket, e error)) {
	p.sio.On("error", fn)
}

func (p *Pico) On(event string, fn func(msg string)){
	p.sio.On(event, fn) 
}

func (p *Pico) Listen(port int) error {
	envPort := os.Getenv("PORT")
	if len(envPort) > 0 {
		pi, err := strconv.Atoi(envPort)
		if err == nil && pi < 65000 {
			port = pi
		}
	}

	p.server = &graceful.Server{
		Timeout: 3 * time.Second,
		Server: &http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: p.mux,
		},
	}
	return p.server.ListenAndServe()

}

func (p *Pico) Production(){
	isDev = false;
}
func (p *Pico) StopOnInt() {
	p.c = make(chan os.Signal, 1)
	signal.Notify(p.c, os.Interrupt)
	signal.Notify(p.c, os.Kill)
	go func() {
		for sig := range p.c {
			fmt.Println(sig.String(), "Shutting Down!")
			
			close(p.c)
			p.Stop()
			fmt.Println(sig.String(), "Done!")
			os.Exit(0)
		}
	}()
}

func (p *Pico) Stop() {
	p.server.Stop(time.Second * 2)
	fmt.Println("Waiting on Stop Channel")
	<-p.server.StopChan()
	fmt.Println("Channel Returned")
}

func New() *Pico {
	isDev = true
	return &Pico{mux: httprouter.New()}
}
