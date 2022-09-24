package rpc

import (
	"github.com/oldbai555/lbtool/rpc/codec"
	"reflect"
)

// request stores all information of a call
type request struct {
	h            *codec.Header // header of request
	argv, replyv reflect.Value // argv and replyv of request
	mtype        *methodType
	svc          *service
}
