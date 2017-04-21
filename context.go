package picoweb

import "net/http"
import "io/ioutil"
import "html/template"
import "fmt"
import "encoding/json"
import "github.com/googollee/go-socket.io"
import "time"
import "strconv"
import "github.com/pkg/errors"

type Context struct {
	w         http.ResponseWriter
	r         *http.Request
	params    map[string]string
	SessionId string
}

func (c *Context) Body() ([]byte, error) {
	bts, err := ioutil.ReadAll(c.r.Body)
	return bts, err
}

func (c *Context) Bytes() []byte {
	bts, _ := ioutil.ReadAll(c.r.Body)
	return bts
}

func (c *Context) Query(key string) string {
	return c.r.URL.Query().Get(key)
}

func (c *Context) QueryInt(key string) (int, error) {
	v := c.Query(key)
	if len(v) == 0 {
		return 0, errors.New("key not found in path")
	}
	i, err := strconv.Atoi(key)
	if err != nil {
		return 0, errors.New("Could not parse as Int")
	}
	return i, nil
}
func (c *Context) QueryBool(key string) (bool, error) {
	v := c.Query(key)
	if len(v) == 0 {
		return false, errors.New("key not found in path")
	}
	b, err := strconv.ParseBool(key)
	if err != nil {
		return false, errors.New("Could not parse as Bool")
	}
	return b, nil
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

func (c *Context) RemoteIP() string {
	return c.r.RemoteAddr
}

func (c *Context) R() *http.Request {
	return c.r
}

func (c *Context) SetHeader(key string, value string) {
	c.w.Header().Set(key, value)
}

func (c *Context) File(filePath string, mimeType string) {
	c.w.Header().Set("content-type", mimeType)
	http.ServeFile(c.w, c.r, filePath)
}

func (c *Context) FileHTML(filePath string) {
	c.w.Header().Set("content-type", "text/html; charset=utf-8")
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

func (c *Context) SetCookie(name, value string, expireIn time.Duration) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   0,
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(expireIn),
		Path:     "/",
		Raw:      value,
		Unparsed: []string{value},
	}
	http.SetCookie(c.w, cookie)
}

func (c *Context) GetCookie(name string) string {
	cookie, err := c.r.Cookie(name)
	if err != nil {
		return ""
	}
	val := cookie.Value
	if len(val) == 0 {
		for _, ck := range c.r.Cookies() {
			if ck.Name == name {
				return ck.Value
			}
		}
	}
	return ""
}

func (c *Context) RemoveCookie(name string) {
	c.SetCookie(name, "", -(time.Hour * 36))
}

type Socket struct {
	socketio.Socket
}
