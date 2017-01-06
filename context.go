package picoweb

import "net/http"
import "io/ioutil"
import "html/template"
import "fmt"
import "encoding/json"

type Context struct {
	w      http.ResponseWriter
	r      *http.Request
	params map[string]string
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

func (c *Context) Method() string {
	return c.r.Method
}

func (c *Context) Header(key string) string {
	return c.r.Header.Get(key)
}

func (c *Context) SetHeader(key string, value string) {
	c.w.Header().Set(key, value)
}

func (c *Context) File(filePath string, mimeType string) {
	c.w.Header().Set("content-type", mimeType)
	http.ServeFile(c.w, c.r, filePath)
}

func (c *Context) String(str string) {
	fmt.Fprint(c, str)
}

func (c *Context) Status(statusCode int) {
	c.w.WriteHeader(statusCode)
}

func (c *Context) Json(data interface{}) {
	jsoned, _ := json.Marshal(data)
	c.ResponseHeader().Add("content-type", "application/json")
	fmt.Fprint(c, string(jsoned))
}

func (c *Context) View(filePath string, data interface{}) {
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		fmt.Fprint(c, err.Error())
		return
	}
	err = tmpl.Execute(c.w, data)
}

func (c *Context) Params(name string) string {
	v, _ := c.params[name]
	return v
}

func (c *Context) ResponseHeader() http.Header {
	return c.w.Header()
}

func (c *Context) WriteHeader(n int) {
	c.w.WriteHeader(n)
}

func (c *Context) Write(b []byte) (int, error) {
	return c.w.Write(b)
}
