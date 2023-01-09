package picoweb

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/julienschmidt/httprouter"
)

type middlewarehandler func(c *Context) bool

var premiddlewares []middlewarehandler
var postmiddlewares []middlewarehandler
var middlewares []middlewarehandler
var must middlewarehandler

func middle(p PicoHandler, appname string, useAppManager bool) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if r := recover(); r != nil {
			fmt.Println("Recovering in Middle")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var sessionId string
		if useAppManager {
			if len(appname) > 0 {
				c, err := r.Cookie(appname)
				if err == nil {
					sessionId = c.Value
				}
			}
		}

		atomic.AddUint64(&requestCount, 1)

		start := time.Now()
		c := &Context{SessionId: sessionId, AppName: appname, UserManager: &usermanager{appname: appname}, w: w, r: r, params: make(map[string]string), Start: time.Now()}

		for _, par := range ps {
			c.params[par.Key] = par.Value
		}

		//w.Header().Set("Access-Control-Allow-Origin", "*")

		docontinue := !skipmiddlewares

		if docontinue {
			for _, m := range premiddlewares {
				docontinue = m(c)
				if !docontinue {
					break
				}
			}
		}

		if docontinue {
			for _, m := range middlewares {
				docontinue = m(c)
				if !docontinue {
					break
				}
			}
		}

		if skipmiddlewares || docontinue {
			p(c)
		}

		if docontinue {
			for _, m := range postmiddlewares {
				docontinue = m(c)
				if !docontinue {
					break
				}
			}
		}

		if must != nil {
			must(c)
		}

		if c.s != nil {
			c.s.Close()
		}

		if c.red != nil {
			c.red.Close()
		}

		if isDev {
			fmt.Println(time.Since(start), r.URL, atomic.LoadUint64(&requestCount))
		}
	}
}

func middlehttp(fn http.Handler) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		//start := time.Now()
		fn.ServeHTTP(w, r)
		//fmt.Println(time.Since(start), r.URL, requestCount)
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
