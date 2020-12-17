package base

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

// HandlerFunc defines the request handler used by webgo.
type HandlerFunc func(*Context)

// Engine implement the interface of HTTP Server
type Engine struct {
	router *router
	*RouterGroup
	groups        []*RouterGroup
	htmlTemplates *template.Template
	// User defined render functino
	funcMap template.FuncMap
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

// Use is defined to add middleware to group.
func (rg *RouterGroup) Use(middlewares ...HandlerFunc) {
	rg.middlewares = append(rg.middlewares, middlewares...)
}

// GET defines the method to add GET request.
func (rg *RouterGroup) GET(pattern string, handler HandlerFunc) {
	rg.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request.
func (rg *RouterGroup) POST(pattern string, handler HandlerFunc) {
	rg.addRoute("POST", pattern, handler)
}

// createStaticHandler receives rendering files' path
// and return rendering handler.
func (rg *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(rg.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// Static server static files for users.
func (rg *RouterGroup) Static(relativePath string, root string) {
	handler := rg.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	rg.GET(urlPattern, handler)
}

// SetFuncMap allow user define render function by themselves.
func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

// LoadHTMLGlob load all render functions.
func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}

// Run a http server.
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

// ServeHTTP is the unified entrance of request.
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		// 判断用户组
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			// 读取该用户组的所有中间件
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	// 加载所有符合条件的中间件到context中
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = e
	e.router.handle(c)
}

// Default use Logger() & Recovery() middlewares by default.
func Default() *Engine {
	e := New()
	e.Use(Logger(), Recovery())
	return e.engine
}
