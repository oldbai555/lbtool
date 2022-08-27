package web

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/oldbai555/lb/comm"
	"net/http"
)

type H map[string]interface{}

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

	//
	ctx       context.Context
	parentCtx *context.Context

	//
	serverName string
	hint       string
}

func (c *Context) GetCtx() context.Context {
	return c.ctx
}

func (c *Context) GetServiceName() string {
	return c.serverName
}

func (c *Context) GetHint() string {
	return c.hint
}

func newContext(w http.ResponseWriter, req *http.Request, ctx context.Context, serverName string) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,

		ctx: ctx,

		serverName: serverName,
		hint:       comm.GetRandomString(16, comm.RandomStringModNumberPlusLetter),
	}
}

// Param 拿到路径上的参数
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// PostForm 解析参数
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 解析参数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status 设置http状态
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 设置请求头
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String 输出字符串
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	if _, err := c.Writer.Write([]byte(fmt.Sprintf(format, values...))); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// JSON json
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data byte流
func (c *Context) Data(code int, data []byte) error {
	c.Status(code)
	_, err := c.Writer.Write(data)
	return err
}

// HTML html网页
func (c *Context) HTML(code int, html string) error {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	_, err := c.Writer.Write([]byte(html))
	return err
}
