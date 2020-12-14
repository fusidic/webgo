package base

import (
	"net/http"
)

// HandlerFunc defines the request handler used by webgo.
type HandlerFunc func(*Context)

// Engine implement the interface of HTTP Server
type Engine struct {
	router *router
}

// New is the constructor of webgo.Engine
func New() *Engine {
	return &Engine{router: newRouter()}
}

func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	e.router.addRouter(method, pattern, handler)
}

// GET defines the method to add GET request.
func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request.
func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.addRoute("POST", pattern, handler)
}

// Run a http server.
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

// ServeHTTP is the unified entrance of request.
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	e.router.handle(c)
}
