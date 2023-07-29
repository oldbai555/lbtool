package event

import (
	"github.com/golang/protobuf/proto"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/routine"
	"sync"
)

type BaseMsg struct {
	IsPb bool // 支持pb类型的数据 或 自定义的数据
	Pb   proto.Message
	Val  interface{}
}

func NewMsg(val interface{}) IMsg {
	return &BaseMsg{Val: val, IsPb: false}
}

func NewPbMsg(pb proto.Message) IMsg {
	return &BaseMsg{Pb: pb, IsPb: true}
}

func (m BaseMsg) GetValue() interface{} {
	if m.IsPb {
		return m.Pb
	}
	return m.Val
}

type BaseEvent struct {
	asyncDoEvent bool
	pool         map[Type]Fns
	lock         sync.RWMutex
}

func NewEvent(ops ...Option) IEvent {
	return buildByOptions()
}

type Option func(*BaseEvent)

func buildByOptions(ops ...Option) *BaseEvent {
	e := &BaseEvent{
		pool: make(map[Type]Fns),
	}
	for i := range ops {
		ops[i](e)
	}
	return e
}

func WithAsyncDoEvent() Option {
	return func(event *BaseEvent) {
		event.asyncDoEvent = true
		return
	}
}

func (e *BaseEvent) Reg(t Type, fn Fn) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.pool[t] = append(e.pool[t], fn)
}

func (e *BaseEvent) Trigger(t Type, m IMsg) {
	e.lock.RLock()
	defer e.lock.RUnlock()
	fns, ok := e.pool[t]
	if !ok {
		log.Warnf("not found type %v", t)
		return
	}
	for _, fn := range fns {
		if e.asyncDoEvent {
			routine.GoV2(func() error {
				err := fn(m)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}
				return nil
			})
			continue
		}
		err := fn(m)
		if err != nil {
			log.Errorf("err:%v", err)
			continue
		}
	}
}
