package routine

import (
	"context"
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/alarm"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/utils"
	"runtime/debug"
	"strings"
)

func Go(ctx context.Context, logic func(ctx context.Context) error) {
	// 可以考虑放 traceId 链路追踪
	go func() {
		defer CatchPanic(func(err interface{}) {
			// 	捕获错误后的补救行为
			alarm.Default("system").Alert("Error", fmt.Sprintf("err:%v\ntime:%d", err, utils.TimeNow()))
		})
		err := logic(ctx)
		if lberr.GetErrCode(err) < 0 {
			msg := fmt.Sprintf("moduleName : go-routine err %v", err)
			log.Errorf(msg)
			// 错误通知
		}
	}()
}

func Run(logic func() error) {
	defer CatchPanic(func(err interface{}) {
		alarm.Default("system").Alert("Error", fmt.Sprintf("err:%v\ntime:%d", err, utils.TimeNow()))
	})
	err := logic()
	if lberr.GetErrCode(err) < 0 {
		msg := fmt.Sprintf("moduleName : go-routine err %v", err)
		log.Errorf(msg)
		// 错误通知
	}
}

func GoV2(fn func() error) {
	go func() {
		defer CatchPanic(func(err interface{}) {
			// 	捕获错误后的补救行为
			alarm.Default("system").Alert("Error", fmt.Sprintf("err:%v\ntime:%d", err, utils.TimeNow()))
		})
		err := fn()
		if lberr.GetErrCode(err) < 0 {
			msg := fmt.Sprintf("moduleName : go-routine err %v", err)
			log.Errorf(msg)
			// 错误通知
		}
	}()
}

func CatchPanic(panicCallback func(err interface{})) {
	if err := recover(); err != any(nil) {
		log.Errorf("PROCESS PANIC: err %s", err)
		st := debug.Stack()
		if len(st) > 0 {
			log.Errorf("dump stack (%s):", err)
			lines := strings.Split(string(st), "\n")
			for _, line := range lines {
				log.Errorf("%s", line)
			}
		} else {
			log.Errorf("stack is empty (%s)", err)
		}
		if panicCallback != nil {
			panicCallback(err)
		}
	}
}
