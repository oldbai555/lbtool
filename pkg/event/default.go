package event

import (
	"github.com/oldbai555/lbtool/log"
	"sync"
)

var defaultEvent IEvent
var once sync.Once

func Reg(t Type, fn Fn) {
	once.Do(func() {
		if defaultEvent == nil {
			defaultEvent = NewBaseEvent()
		}
	})
	defaultEvent.Reg(t, fn)
}

func Trigger(t Type, m IMsg) {
	if defaultEvent == nil {
		log.Warnf("not found default event , type is %v", t)
		return
	}
	defaultEvent.Trigger(t, m)
}
