package log

import (
	"context"
	"errors"
	"fmt"
	fmt2 "github.com/oldbai555/lb/log/fmt"
	"github.com/oldbai555/lb/log/writer"
	"sync"
)

const (
	PROD = "PROD"
	DEV  = "DEV"
	TEST = "TEST"
	DEMO = "DEMO"
)

// Logger 日志业务
type Logger struct {
	logLevel     fmt2.Level
	loggerWriter writer.LoggerWriter
	mu           sync.RWMutex
}

func NewDefaultLogger(e string) *Logger {
	return &Logger{
		loggerWriter: writer.NewDefaultSimpleLoggerWriter(e),
	}
}

func (l *Logger) SetLogLevel(level fmt2.Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logLevel = level
}

func (l *Logger) GetLogLevel() fmt2.Level {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.logLevel
}

func (l *Logger) Debugf(ctx context.Context, format string, args ...interface{}) {
	if err := l.write(ctx, fmt2.LevelDebug, append([]interface{}{format}, args...)...); err != nil {
		panic(any(err))
	}
}

func (l *Logger) Infof(ctx context.Context, format string, args ...interface{}) {
	if err := l.write(ctx, fmt2.LevelInfo, append([]interface{}{format}, args...)...); err != nil {
		panic(any(err))
	}
}

func (l *Logger) Warnf(ctx context.Context, format string, args ...interface{}) {
	if err := l.write(ctx, fmt2.LevelWarn, append([]interface{}{format}, args...)...); err != nil {
		panic(any(err))
	}
}

func (l *Logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	if err := l.write(ctx, fmt2.LevelError, append([]interface{}{format}, args...)...); err != nil {
		panic(any(err))
	}
}

func (l *Logger) write(ctx context.Context, level fmt2.Level, args ...interface{}) error {
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

	if err := l.loggerWriter.Write(ctx, level, format, realArgs...); err != nil {
		return err
	}

	return nil
}

func (l *Logger) Flush() error {
	return l.loggerWriter.Flush()
}
