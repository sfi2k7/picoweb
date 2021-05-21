package main

import (
	"fmt"

	"github.com/sfi2k7/picoweb"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	p := picoweb.New()
	p.StopOnInt()
	p.Get("/", func(c *picoweb.Context) {

		s, err := c.Mongo()
		if err != nil {
			fmt.Println(err)
			return
		}
		defer s.Close()
		var all []bson.M
		s.DB("local").C("startup_log").Find(bson.M{}).All(&all)

		c.View("./index.html", struct {
			Name string
			All  interface{}
		}{
			All:  all,
			Name: "Faisal",
		})
	})
	fmt.Println(p.Listen(57432))
}
