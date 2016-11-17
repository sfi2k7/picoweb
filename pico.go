package picoweb

import (
	"net/http"
	"os"
	"time"

	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tylerb/graceful"
)

var (
	RequestCount int
)

type Pico struct {
	mux    *httprouter.Router
	server *graceful.Server
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

func (p *Pico) Listen(port int) error {
	envPort := os.Getenv("PORT")
	if len(envPort) > 0 {
		pi, err := strconv.Atoi(envPort)
		if err == nil && pi < 65000 {
			port = pi
		}
	}

	p.server = &graceful.Server{
		Timeout: 10 * time.Second,
		Server: &http.Server{
			Addr:    ":" + os.Getenv("PORT"),
			Handler: p.mux,
		},
	}
	return p.server.ListenAndServe()

}

func (p *Pico) Stop() {
	p.server.Stop(time.Second * 5)
	<-p.server.StopChan()
}

func New() *Pico {
	return &Pico{mux: httprouter.New()}
}
