package rpc

type ServerCmdList []ServerCmd

type IFunc func(ctx *Context, req *IRequest) (*IResponse, error)

type ServerCmd struct {
	ServerName string      `json:"serverName"`
	Path       string      `json:"path"`
	IFun       interface{} `json:"i_fun"`
}

type Context struct {
}

type IRequest interface {
	GetHeader(name string) (string, error)
}

type IResponse interface {
	SetHeader(name string, value string) error
}
