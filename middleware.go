package picoweb

import (
	"fmt"
	"net/http"
	"time"

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
