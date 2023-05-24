package main

import (
	"github.com/oldbai555/lbtool/demo/rpc_client"
	rpc2 "github.com/oldbai555/lbtool/discard/rpc"
	"github.com/oldbai555/lbtool/log"
)

func main() {
	var foo rpc_client.Foo
	if err := rpc2.Register(&foo); err != nil {
		log.Errorf("register error:%v", err)
		return
	}
	err := rpc2.ServerRun(9999)
	if err != nil {
		log.Errorf("err is :%v", err)
		return
	}
}

var serverCmdList = []rpc2.ServerCmd{
	{
		ServerName: "server",
		IFun:       Sum,
	},
}

type FooReq struct {
}

func (f *FooReq) GetHeader(name string) (string, error) {
	return "", nil
}

var _ rpc2.IRequest = (*FooReq)(nil)

type FooRsp struct {
}

func (f *FooRsp) SetHeader(name string, value string) error {
	return nil
}

var _ rpc2.IResponse = (*FooRsp)(nil)

func Sum(ctx *rpc2.Context, req *FooReq) (*FooRsp, error) {
	var rsp FooRsp
	return &rsp, nil
}
