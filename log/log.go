package log

import (
	"errors"
	"fmt"
	"github.com/oldbai555/lbtool/log/iface"
	"github.com/oldbai555/lbtool/utils"
	"github.com/petermattis/goid"
	"io"
	"sync"
)

var (
	log        *logger
	logCtx     = map[int64]string{}
	logCtxMu   sync.RWMutex
	moduleName = "UNKNOWN"
)

func init() {
	log = newLogger()
}

func SetLogLevel(level utils.Level) {
	if log == nil {
		log = newLogger()
	}
	log.logLevel = level
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

func SetModuleName(name string) {
	moduleName = name
}

func GetWriter() io.Writer {
	return log.logWriter
}

func GetLogger() *logger {
	return log
}

func Debugf(format string, args ...interface{}) {

	if err := log.write(utils.LevelDebug, append([]interface{}{format}, args...)...); err != nil {
		panic(any(err))
	}
}

func Infof(format string, args ...interface{}) {
	if err := log.write(utils.LevelInfo, append([]interface{}{format}, args...)...); err != nil {
		panic(any(err))
	}
}

func Warnf(format string, args ...interface{}) {
	if err := log.write(utils.LevelWarn, append([]interface{}{format}, args...)...); err != nil {
		panic(any(err))
	}

}

func Errorf(format string, args ...interface{}) {
	if err := log.write(utils.LevelError, append([]interface{}{format}, args...)...); err != nil {
		panic(any(err))
	}
}

//===================================logger===================================================

// Logger 日志业务
type logger struct {
	logLevel  utils.Level
	logWriter iface.LogWriter
	fmt       iface.Formatter
	mu        sync.RWMutex
}

func newLogger() *logger {
	return &logger{
		logWriter: newLogWriterImpl(),
		fmt:       newSimpleFormatter(),
	}
}

func (l *logger) SetSkipCall(skipCall int) {
	l.fmt.SetSkipCall(skipCall)
}

func (l *logger) write(level utils.Level, args ...interface{}) error {
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

	stdoutColor := utils.LevelToStdoutColorMap[level]
	logContent, err := l.fmt.Sprintf(level, stdoutColor, fmt.Sprintf(format, realArgs...))
	if err != nil {
		return err
	}

	if _, err := l.logWriter.Write([]byte(logContent)); err != nil {
		return err
	}

	return nil
}

func (l *logger) Flush() error {
	return l.logWriter.Flush()
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *logger) Printf(format string, v ...any) {
	if err := log.write(utils.LevelInfo, append([]interface{}{format}, v...)...); err != nil {
		panic(any(err))
	}

}
