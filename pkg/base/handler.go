package base

import (
	"log"
	"net/http"
)

// HandlerFunc defines the request handler used by webgo.
type HandlerFunc func(*Context)

// Engine implement the interface of HTTP Server
type Engine struct {
	router *router
	*RouterGroup
	groups []*RouterGroup
}

// RouterGroup controls those middlewares' implement by groups.
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	engine      *Engine // All groups share one Engine.
}

// New is the constructor of webgo.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup} // 按理说这个不该加入进来的，逻辑层次上有些混淆
	return engine
}

// Group is defined to create a new RouterGroup
func (rg *RouterGroup) Group(prefix string) *RouterGroup {
	engine := rg.engine
	newGroup := &RouterGroup{
		prefix: rg.prefix + prefix,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (rg *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := rg.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	rg.engine.router.addRouter(method, pattern, handler)
}

// GET defines the method to add GET request.
func (rg *RouterGroup) GET(pattern string, handler HandlerFunc) {
	rg.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request.
func (rg *RouterGroup) POST(pattern string, handler HandlerFunc) {
	rg.addRoute("POST", pattern, handler)
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
