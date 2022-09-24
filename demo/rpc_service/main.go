package main

import (
	"github.com/oldbai555/lbtool/demo/rpc_client"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/rpc"
)

func main() {
	var foo rpc_client.Foo
	if err := rpc.Register(&foo); err != nil {
		log.Errorf("register error:%v", err)
		return
	}
	err := rpc.ServerRun(9999)
	if err != nil {
		log.Errorf("err is :%v", err)
		return
	}
}

var serverCmdList = []rpc.ServerCmd{
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

var _ rpc.IRequest = (*FooReq)(nil)

type FooRsp struct {
}

func (f *FooRsp) SetHeader(name string, value string) error {
	return nil
}

var _ rpc.IResponse = (*FooRsp)(nil)

func Sum(ctx *rpc.Context, req *FooReq) (*FooRsp, error) {
	var rsp FooRsp
	return &rsp, nil
}
