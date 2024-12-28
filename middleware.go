package picoweb

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/julienschmidt/httprouter"
)

type PicoMiddleWareHandler func(c *Context) bool

var premiddlewares []PicoMiddleWareHandler
var postmiddlewares []PicoMiddleWareHandler
var middlewares []PicoMiddleWareHandler
var must []PicoMiddleWareHandler

type reqcount struct {
	r map[string]uint64
	s sync.Mutex
}

func (r *reqcount) Add(k string) {
	r.s.Lock()
	defer r.s.Unlock()

	if r.r == nil {
		r.r = make(map[string]uint64)
	}
	r.r[k]++
}

func (r *reqcount) Get(k string) uint64 {
	r.s.Lock()
	defer r.s.Unlock()

	if r.r == nil {
		return 0
	}
	return r.r[k]
}

var rqc reqcount

func middle(p PicoHandler, appname string) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if r := recover(); r != nil {
			// fmt.Println("Recovering in Middle")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		rqc.Add(r.URL.Path)

		var sessionId string
		// if useAppManager {
		// 	if len(appname) > 0 {
		// 		c, err := r.Cookie(appname)
		// 		if err == nil {
		// 			sessionId = c.Value
		// 		}
		// 	}
		// }

		atomic.AddUint64(&requestCount, 1)

		start := time.Now()
		//UserManager: &usermanager{appname: appname},
		c := &Context{Store: newStore(), SessionId: sessionId, AppName: appname, w: w, r: r, params: make(map[string]string), Start: time.Now()}

		for _, par := range ps {
			c.params[par.Key] = par.Value
		}

		//w.Header().Set("Access-Control-Allow-Origin", "*")

		runNext := !skipmiddlewares
		if runNext {
			for _, m := range premiddlewares {
				runNext = m(c)
				if !runNext {
					break
				}
			}
		}

		if runNext {
			for _, m := range middlewares {
				runNext = m(c)
				if !runNext {
					break
				}
			}
		}

		if runNext || skipmiddlewares {
			p(c)
		}

		if runNext && !skipmiddlewares {
			for _, m := range postmiddlewares {
				runNext = m(c)
				if !runNext {
					break
				}
			}
		}

		// if len(must) > 0 && !skipmiddlewares {
		for _, m := range must {
			runNext = m(c)
			if !runNext {
				break
			}
		}
		// }

		c.params = nil
		c.SessionId = ""
		c.State = nil
		c.Store = nil
		c.User = nil
		c.UserData = nil

		// if c.s != nil {
		// 	c.s.Close()
		// }

		// if c.red != nil {
		// 	c.red.Close()
		// }

		if isDev {
			// fmt.Println("ts:", time.Now().Format(time.RFC3339), time.Since(start), "req#:", rqc.Get(r.URL.Path), "/", atomic.LoadUint64(&requestCount), "url:", r.URL)
			fmt.Printf("ts: %s, time:%s, req:%d/%d, url:%s\n", time.Now().Format(time.RFC1123), time.Since(start), rqc.Get(r.URL.Path), atomic.LoadUint64(&requestCount), r.URL)
		}
	}
}

// func middlehttp(fn http.Handler) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 		//start := time.Now()
// 		fn.ServeHTTP(w, r)
// 		//fmt.Println(time.Since(start), r.URL, requestCount)
// 	}
// }

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
