package picoweb

import (
	"fmt"
	"net/http"
	"time"

	"github.com/googollee/go-socket.io"
	"github.com/julienschmidt/httprouter"
)

func middle(p PicoHandler) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		RequestCount++
		start := time.Now()
		c := &Context{w: w, r: r, params: make(map[string]string)}

		for _, par := range ps {
			c.params[par.Key] = par.Value
		}

		//w.Header().Set("Access-Control-Allow-Origin", "*")

		p(c)

		if c.s != nil {
			c.s.Close()
		}

		if c.red != nil {
			c.red.Close()
		}

		if isDev {
			fmt.Println(time.Since(start), r.URL, RequestCount)
		}
	}
}

func middlehttp(fn http.Handler) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		RequestCount++
		//start := time.Now()
		fn.ServeHTTP(w, r)
		//fmt.Println(time.Since(start), r.URL, RequestCount)
	}
}

func socket_middle(fn func(c Socket)) func(socketio.Socket) {
	return func(socket socketio.Socket) {
		cs := Socket{Socket: socket}
		fn(cs)
	}
}

func socket_middle_error(fn func(c Socket, err error)) func(socketio.Socket, error) {
	return func(socket socketio.Socket, err error) {
		cs := Socket{Socket: socket}
		fn(cs, err)
	}
}
