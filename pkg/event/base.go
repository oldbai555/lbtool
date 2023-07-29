package event

import (
	"github.com/golang/protobuf/proto"
	"github.com/oldbai555/lbtool/log"
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
	pool map[Type]Fns
	lock sync.RWMutex
}

func NewBaseEvent() IEvent {
	return &BaseEvent{
		pool: make(map[Type]Fns),
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
		err := fn(m)
		if err != nil {
			log.Errorf("err:%v", err)
			continue
		}
	}
}
