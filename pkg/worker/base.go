package worker

import (
	"context"
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/pkg/routine"
	"github.com/oldbai555/lbtool/pkg/signal"
	"sync"
)

var _ IWorker = (*BaseWorker)(nil)

const defaultItemMax = 1024

type BaseWorker struct {
	c    chan IMsg
	svr  string
	size int

	sync.Mutex
}

func NewWorker(size int, svr string) IWorker {
	return &BaseWorker{
		c:    make(chan IMsg, size),
		svr:  svr,
		size: size,
	}
}

func (e *BaseWorker) Start(fn func(msg IMsg)) {
	routine.GoV2(func() (err error) {
		log.Infof(fmt.Sprintf("============ %s starting service ============", e.svr))
		defer func() {
			log.Infof(fmt.Sprintf("============ %s end service ============", e.svr))
		}()
		var msgs []IMsg
		for {
			select {
			case <-signal.GetSignalChan():
				log.Infof(fmt.Sprintf("============ %s service done ============", e.svr))
				e.Stop()
				return
			case receive := <-e.c:
				msgs = append(msgs, receive)
				msgs = append(msgs, e.receive()...)
				for i := 0; i < len(msgs); i++ {
					log.Infof("msg[%d] is %v", i, msgs[i].GetValue())
					fn(msgs[i])
				}
				msgs = nil
			}
		}
	})
}

func (e *BaseWorker) Send(t Type, v interface{}) (err error) {
	//if e.c == nil {
	//	e.c = make(chan IMsg, e.size)
	//}
	select {
	case e.c <- NewMsg(t, v):
		log.Infof("send msg type is %v", t)
	default: // 不让消息丢失的话 需要做兜底
		log.Warnf("%s worker queue is full", e.svr)
		err = lberr.NewErr(1002, fmt.Sprintf("%s worker queue is full", e.svr))
	}
	return
}

func (e *BaseWorker) Stop() {
	close(e.c)
}

func (e *BaseWorker) receive() []IMsg {
	var vs []IMsg
OUT:
	for {
		if len(vs) >= defaultItemMax {
			break
		}
		select {
		case receive := <-e.c:
			vs = append(vs, receive)
		default:
			log.Debugf("not receive srv is %v ...,", e.svr)
			break OUT
		}
	}
	return vs
}

// =================================================================

var _ IHandlerMgr = (*BaseHandlerMgr)(nil)

type BaseHandlerMgr struct {
	fnM map[Type]DoHandlerFn
	sync.Mutex
}

func NewHandlerMgr() *BaseHandlerMgr {
	return &BaseHandlerMgr{
		fnM: make(map[Type]DoHandlerFn),
	}
}

func (e *BaseHandlerMgr) Register(t Type, fn DoHandlerFn) {
	e.Lock()
	defer e.Unlock()
	_, ok := e.fnM[t]
	if ok {
		panic(fmt.Sprintf("already registered fn , type is %v", t))
		return
	}
	e.fnM[t] = fn
	return
}

func (e *BaseHandlerMgr) Call(ctx context.Context, m IMsg) error {
	fn, ok := e.fnM[m.GetType()]
	if !ok {
		log.Warnf("not found fn , type is %v", m.GetType())
		return lberr.NewInvalidArg("not found fn , type is %v", m.GetType())
	}
	err := fn(ctx, m.GetType())
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

// =================================================================

var _ IMsg = (*BaseMsg)(nil)

type BaseMsg struct {
	Typ Type
	Val interface{}
}

func NewMsg(typ Type, val interface{}) IMsg {
	return &BaseMsg{Typ: typ, Val: val}
}

func (i *BaseMsg) GetType() Type {
	return i.Typ
}

func (i *BaseMsg) GetValue() interface{} {
	return i.Val
}
