package gormx

import (
	"context"
	"errors"
	"fmt"
	"github.com/oldbai555/gorm/logger"
	"github.com/oldbai555/gorm/utils"
	"github.com/oldbai555/lbtool/log"
	"time"
)

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

const (
	traceStr     = "%s [%.3fms] [rows:%v] %s"
	traceWarnStr = "%s %s [%.3fms] [rows:%v] %s"
	traceErrStr  = "%s %s [%.3fms] [rows:%v] %s"
)

// NewOrmLog initialize logger
func NewOrmLog(slowThreshold time.Duration) logger.Interface {
	return &ormlog{
		slowThreshold: slowThreshold,
	}
}

type ormlog struct {
	slowThreshold time.Duration
}

func (l *ormlog) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l ormlog) Info(ctx context.Context, msg string, data ...interface{}) {
	log.Infof(msg, data...)
}

func (l ormlog) Warn(ctx context.Context, msg string, data ...interface{}) {
	log.Warnf(msg, data...)

}

func (l ormlog) Error(ctx context.Context, msg string, data ...interface{}) {
	log.Errorf(msg, data...)
}

func (l ormlog) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil && (!errors.Is(err, logger.ErrRecordNotFound)):
		sql, rows := fc()
		if rows == -1 {
			l.Error(ctx, traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Error(ctx, traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.slowThreshold && l.slowThreshold != 0:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.slowThreshold)
		if rows == -1 {
			l.Warn(ctx, traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Warn(ctx, traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	default:
		sql, rows := fc()
		if rows == -1 {
			l.Info(ctx, traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Info(ctx, traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
