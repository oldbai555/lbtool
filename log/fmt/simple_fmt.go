package fmt

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"
)

const (
	DefaultSkipCall   = 5
	DefaultFormatType = FormatText
)

// SimpleFormatter 格式化日志
type SimpleFormatter struct {
	skipCall   int    // 用于调用层数
	formatType Format // 格式化类型
}

func NewDefaultSimpleFormatter() *SimpleFormatter {
	return &SimpleFormatter{
		skipCall:   DefaultSkipCall,
		formatType: FormatText,
	}
}

func (s *SimpleFormatter) Sprintf(ctx context.Context, level Level, color Color, format string, args ...interface{}) (string, error) {
	// 日志等级
	levelStr, err := TransferLevelToStr(level)
	if err != nil {
		return "", err
	}

	// 字体颜色
	colorStdout, err := GetColorStdout(color)
	if err != nil {
		return "", err
	}

	// skip是层数，调用Caller函数外层的函数。1代表上次，2代表上上层，一般我们需要定位的也就是行数line跟file文件名
	pc, callFile, callLine, ok := runtime.Caller(s.skipCall)
	var callFuncName string
	if ok {
		// 拿到调用方法
		callFuncName = runtime.FuncForPC(pc).Name()
	}

	now := time.Now()
	rawFormatted := fmt.Sprintf(format, args...)

	switch s.formatType {
	case FormatText:
		return fmt.Sprintf(
			"%s-%s.%04d %s %s:%d:%s ==> %s \u001B[0m\n",
			colorStdout,
			now.Format("2006-01-02T15:04:05"),
			now.Nanosecond()/100000,
			levelStr,
			callFile,
			callLine,
			callFuncName,
			rawFormatted,
		), nil
	default:
		return "", errors.New("not support log format")
	}
}

var _ Formatter = (*SimpleFormatter)(nil)
