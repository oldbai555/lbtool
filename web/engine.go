package web

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

const (
	GET  = "GET"
	POST = "POST"
)

// HandlerFunc defines the request handler used by web
type HandlerFunc func(c *Context) error

// Engine implement the internal of ServeHTTP
type Engine struct {
	env        string
	serverName string
	port       uint32

	*RouterGroup
	router *router
	groups []*RouterGroup // store all groups

	//htmlTemplates *template.Template // 将所有的模板加载进内存
	//funcMap       template.FuncMap   // 所有的自定义模板渲染函数
}

var _ http.Handler = (*Engine)(nil)

// New is the constructor of gee.Engine
func New(serverName string, port uint32) *Engine {
	e := &Engine{
		serverName: serverName,
		port:       port,

		router: newRouter(),
	}
	e.RouterGroup = &RouterGroup{engine: e}
	e.groups = []*RouterGroup{e.RouterGroup}
	e.Use(Recovery())
	return e
}

//
//func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
//	engine.funcMap = funcMap
//}
//
//func (engine *Engine) LoadHTMLGlob(pattern string) {
//	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
//}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute(GET, pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute(POST, pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run() (err error) {
	return http.ListenAndServe(fmt.Sprintf(":%d", engine.port), engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := newContext(w, req, context.TODO(), engine.serverName)
	c.handlers = middlewares
	//c.engine = engine 先不支持HTML渲染
	engine.router.handle(c)
}
