package picoweb

import (
	"fmt"
	"net/http"
	"time"
)

func middle(p PicoHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		c := &Context{}
		c.r = r
		c.w = w
		p(c)
		fmt.Println("Took", time.Since(start))
	}
}
