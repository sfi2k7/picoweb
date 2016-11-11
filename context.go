package picoweb

import "net/http"
import "io/ioutil"
import "html/template"
import "fmt"
import "encoding/json"

type Context struct {
	w http.ResponseWriter
	r *http.Request
}

func (c *Context) Body() ([]byte, error) {
	bts, err := ioutil.ReadAll(c.r.Body)
	return bts, err
}

func (c *Context) Query(key string) string {
	return c.r.URL.Query().Get(key)
}

func (c *Context) Form(key string) string {
	return c.r.FormValue(key)
}

func (c *Context) Json(data interface{}) {
	jsoned, _ := json.Marshal(data)
	c.Header().Add("content-type", "application/json")
	fmt.Fprint(c, string(jsoned))
}

func (c *Context) View(filePath string, data interface{}) {
	tmpl, err := template.New("temp").ParseFiles(filePath)
	if err != nil {
		fmt.Fprint(c, err.Error())
		return
	}
	tmpl.Execute(c, data)
}

func (c *Context) Header() http.Header {
	return c.w.Header()
}

func (c *Context) WriteHeader(n int) {
	c.w.WriteHeader(n)
}

func (c *Context) Write(b []byte) (int, error) {
	return c.w.Write(b)
}
