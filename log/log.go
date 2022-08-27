package log

import (
	"errors"
	"fmt"
	"github.com/oldbai555/lb/comm"
	"github.com/petermattis/goid"
	"sync"
)

var (
	log        *logger
	logCtx     = map[int64]string{}
	logCtxMu   sync.RWMutex
	env        string
	moduleName = "UNKNOWN"
)

func init() {
	if env == "" {
		env = comm.DEV
	}
	log = newLogger(env)
}

func SetLogHint(hint string) {
	i := goid.Get()
	logCtxMu.Lock()
	if hint == "" {
		delete(logCtx, i)
	} else {
		logCtx[i] = hint
	}
	logCtxMu.Unlock()
}

func getLogHint() string {
	i := goid.Get()
	logCtxMu.RLock()
	v := logCtx[i]
	logCtxMu.RUnlock()
	return v
}

func SetEnv(e string) {
	env = e
}

func SetModuleName(name string) {
	moduleName = name
}

func Debugf(format string, args ...interface{}) {

	if err := log.write(levelDebug, append([]interface{}{format}, args...)...); err != nil {
		panic(any(err))
	}
}

func Infof(format string, args ...interface{}) {
	if err := log.write(levelInfo, append([]interface{}{format}, args...)...); err != nil {
		panic(any(err))
	}
}

func Warnf(format string, args ...interface{}) {
	if err := log.write(levelWarn, append([]interface{}{format}, args...)...); err != nil {
		panic(any(err))
	}

}

func Errorf(format string, args ...interface{}) {
	if err := log.write(levelError, append([]interface{}{format}, args...)...); err != nil {
		panic(any(err))
	}
}

//===================================logger===================================================

// Logger 日志业务
type logger struct {
	logLevel  Level
	logWriter logWriter
	mu        sync.RWMutex
}

func newLogger(e string) *logger {
	return &logger{
		logWriter: newLogWriterImpl(e),
	}
}

func (l *logger) write(level Level, args ...interface{}) error {
	if l.logLevel > level {
		return nil
	}

	argNum := len(args)
	if argNum == 0 {
		return errors.New("args num is 0")
	}

	var realArgs []interface{}
	if argNum > 1 {
		realArgs = args[1:]
	}

	var (
		format string
		ok     bool
	)
	if format, ok = args[0].(string); !ok {
		format = fmt.Sprint(format)
	}

	if err := l.logWriter.Write(level, fmt.Sprintf(format, realArgs...)); err != nil {
		return err
	}

	return nil
}

func (l *logger) Flush() error {
	return l.logWriter.Flush()
}
