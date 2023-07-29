package worker

import "context"

type IHandlerMgr interface {
	Register(Type, DoHandlerFn)
	Call(context.Context, IMsg) error
}

type Type uint32
type DoHandlerFn func(context.Context, interface{}) error
