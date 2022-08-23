package web

import (
	"fmt"
	"net/http"
)

const (
	GET  = "GET"
	POST = "POST"
)

// HandlerFunc defines the request handler used by web
type HandlerFunc func(c *Context) error

// Engine implement the interface of ServeHTTP
type Engine struct {
	router *router
}

var _ http.Handler = (*Engine)(nil)

// New is the constructor of gee.Engine
func New() *Engine {
	return &Engine{router: newRouter()}
}

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
func (engine *Engine) Run(port uint32) (err error) {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
