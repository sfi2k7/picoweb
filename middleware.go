package picoweb

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/julienschmidt/httprouter"
)

func middle(p PicoHandler) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if r := recover(); r != nil {
			fmt.Println("Recovering in Middle")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		atomic.AddUint64(&RequestCount, 1)

		start := time.Now()
		c := &Context{w: w, r: r, params: make(map[string]string), Start: time.Now()}

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
			fmt.Println(time.Since(start), r.URL, atomic.LoadUint64(&RequestCount))
		}
	}
}

func middlehttp(fn http.Handler) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		//start := time.Now()
		fn.ServeHTTP(w, r)
		//fmt.Println(time.Since(start), r.URL, RequestCount)
	}
}

// func socket_middle(fn func(c Socket)) func(socketio.Socket) {
// 	return func(socket socketio.Socket) {
// 		cs := Socket{Socket: socket}
// 		fn(cs)
// 	}
// }

// func socket_middle_error(fn func(c Socket, err error)) func(socketio.Socket, error) {
// 	return func(socket socketio.Socket, err error) {
// 		cs := Socket{Socket: socket}
// 		fn(cs, err)
// 	}
// }
