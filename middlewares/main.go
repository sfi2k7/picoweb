package middlewares

import "github.com/sfi2k7/picoweb"

var TestMiddleware = func(c *picoweb.Context) bool {
	return true
}
