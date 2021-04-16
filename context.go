package picoweb

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/redis.v4"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	mgo "gopkg.in/mgo.v2"
)

type Context struct {
	w         http.ResponseWriter
	r         *http.Request
	params    map[string]string
	SessionId string
	s         *mgo.Session
	red       *redis.Client
	machineID string
	Start     time.Time
}

func (c *Context) SessionHash() string {
	hasher := sha1.New()
	hasher.Write([]byte(c.r.UserAgent()))
	hasher.Write([]byte(c.r.RemoteAddr))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
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

func (c *Context) QueryCaseIn(key string) string {
	for k, v := range c.r.URL.Query() {
		if strings.ToLower(k) == strings.ToLower(key) {
			if len(v) > 0 {
				return v[0]
			}
			return ""
		}
	}
	return ""
}

func (c *Context) QueryInt(key string) (int, error) {
	v := c.Query(key)
	if len(v) == 0 {
		return 0, errors.New("key not found in path")
	}
	i, err := strconv.Atoi(v)
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

func (c *Context) URLHashPart() string {
	url := c.r.URL.Path
	i := strings.LastIndex(url, "#")
	if i > 0 {
		return url[i:]
	}
	return ""
}

func (c *Context) BasicAuth() (string, string, bool) {
	return c.r.BasicAuth()
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

func (c *Context) Status404() {
	c.w.WriteHeader(http.StatusNotFound)
}

func (c *Context) Status403() {
	c.w.WriteHeader(http.StatusForbidden)
}

func (c *Context) Status401() {
	c.w.WriteHeader(http.StatusUnauthorized)
}

func (c *Context) StatusServerError() {
	c.w.WriteHeader(http.StatusInternalServerError)
}

func (c *Context) Json(data interface{}) (int, error) {
	jsoned, _ := json.Marshal(data)
	c.ResponseHeader().Add("content-type", "application/json")
	return fmt.Fprint(c, string(jsoned))
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

func (c Context) WriteHeader(n int) {
	c.w.WriteHeader(n)
}

func (c Context) Write(b []byte) (int, error) {
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

func (c *Context) URL() *url.URL {
	return c.r.URL
}

func (c *Context) HasPrefix(prefix string) bool {
	return strings.Index(c.r.URL.Path, prefix) == 0
}

func (c *Context) IsStatic() bool {
	p := c.r.URL.Path
	lastSlash := strings.LastIndex(p, "/")

	if lastSlash < 1 {
		return false
	}

	fielName := p[lastSlash:]
	return strings.Index(fielName, ".") > 0
}

func (c *Context) GetStaticFileExt() string {
	return path.Ext(c.r.URL.Path)
}

func (c *Context) Host() string {
	return c.r.Host
}

func (c *Context) Path() string {
	return c.r.URL.Path
}

func (c *Context) W() http.ResponseWriter {
	return c.w
}

func (c *Context) GetStaticDirFile() (string, string) {
	p := c.r.URL.Path
	dir, file := filepath.Split(p)
	return dir, file
}

func (c *Context) GetStaticFile() string {
	_, file := c.GetStaticDirFile()
	return file
}

func (c *Context) GetStaticFilePath() string {
	dir, _ := c.GetStaticDirFile()
	return dir
}

func (c *Context) GetCookie(name string) string {
	cookie, err := c.r.Cookie(name)
	if err != nil {
		fmt.Println("COOKIE ERROR", err)
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
	return val
}

func (c *Context) Mongo() (*mgo.Session, error) {
	if c.s != nil {
		return c.s, nil
	}
	s, err := getSession()
	c.s = s
	return s, err
}

func (c *Context) Redis() (*redis.Client, error) {
	if c.red != nil {
		return c.red, nil
	}

	c.red = redis.NewClient(&redis.Options{
		Addr:     redisURL,
		DB:       0,
		Network:  "tcp",
		Password: redisPassword,
	})
	return c.red, c.red.Ping().Err()
}

func (c *Context) Upgrade() (*websocket.Conn, error) {
	var upgrader = websocket.Upgrader{EnableCompression: true, HandshakeTimeout: time.Second * 5, ReadBufferSize: 4096, WriteBufferSize: 4096}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(c.w, c.r, nil)
	return conn, err
}

func (c *Context) RemoveCookie(name string) {
	c.SetCookie(name, "", -(time.Hour * 36))
}

func (c *Context) Redirect(url string, code int) {
	http.Redirect(c.w, c.r, url, code)
}

// type Socket struct {
// 	socketio.Socket
// }
