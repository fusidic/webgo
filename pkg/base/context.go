package base

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H ...
type H map[string]interface{}

// Context contains the information Request needed.
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int
	// engine pointer
	engine *Engine
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// Next call all defined handlers.
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// Fail only for test
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

// Param ...
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// PostForm ...
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query ...
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status contains the status code.
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader set responses' header. With the form as
// "Content-type": "text/plain" or "application/json" or "text/html"
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String pass string to the response packet.
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON pass JSON to the response packet.
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data pass []byte to packet.
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML pass html string.
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
	// c.Writer.Write([]byte(html))
}
