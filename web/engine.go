package web

import (
	"context"
	"fmt"
	"github.com/oldbai555/lb/log"
	"net/http"
)

const (
	GET  = "GET"
	POST = "POST"
)

// HandlerFunc defines the request handler used by web
type HandlerFunc func(c *Context) error

// Engine implement the internal of ServeHTTP
type Engine struct {
	serverName string
	port       uint32

	*RouterGroup
	router *router
	groups []*RouterGroup // store all groups
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
	return e
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
func (engine *Engine) Run() (err error) {
	return http.ListenAndServe(fmt.Sprintf(":%d", engine.port), engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req, context.TODO(), engine.serverName)
	log.SetLogHint(c.hint)
	engine.router.handle(c)
}
