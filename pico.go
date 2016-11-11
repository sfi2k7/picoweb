package picoweb

import (
	"net/http"
	"os"
	"time"

	"strconv"

	"github.com/tylerb/graceful"
)

type Pico struct {
	mux    *http.ServeMux
	server *graceful.Server
}

type PicoHandler func(c *Context)

func (p *Pico) HandleFunc(pattern string, fn PicoHandler) {
	p.mux.HandleFunc(pattern, middle(fn))
}

func (p *Pico) Static(urlPath, diskPath string) {
	p.mux.Handle(urlPath, http.StripPrefix(urlPath, http.FileServer(http.Dir(diskPath))))
}

func (p *Pico) Listen(port int) {
	os.Setenv("PORT", strconv.Itoa(port))
	p.server = &graceful.Server{
		Timeout: 10 * time.Second,
		Server: &http.Server{
			Addr:    ":" + os.Getenv("PORT"),
			Handler: p.mux,
		},
	}

}

func (p *Pico) Stop() {
	p.server.Stop(time.Second * 5)
	<-p.server.StopChan()
}

func New() *Pico {
	return &Pico{mux: http.NewServeMux()}
}
