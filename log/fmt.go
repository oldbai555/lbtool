package log

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/oldbai555/lbtool/log/iface"
	"github.com/oldbai555/lbtool/utils"
	"github.com/petermattis/goid"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

const (
	DefaultSkipCall = 3
)

var _ iface.Formatter = (*simpleFormatter)(nil)

// SimpleFormatter 格式化日志
type simpleFormatter struct {
	skipCall   int          // 用于调用层数
	formatType utils.Format // 格式化类型
}

func newSimpleFormatter() *simpleFormatter {
	return &simpleFormatter{
		skipCall:   DefaultSkipCall,
		formatType: utils.FormatText,
	}
}

func getPackageName(f string) (filePath string, fileFunc string) {
	slashIndex := strings.LastIndex(f, "/")
	filePath = f
	if slashIndex > 0 {
		idx := strings.Index(f[slashIndex:], ".") + slashIndex
		filePath = f[:idx]
		fileFunc = f[idx+1:]
		return
	}
	return
}

func (s *simpleFormatter) Sprintf(level utils.Level, color utils.Color, buf string) (string, error) {
	// Go获取当前协程信息 第三方库
	var b bytes.Buffer

	routineId := goid.Get()
	pid := os.Getpid()
	// 进程、协程
	b.WriteString(fmt.Sprintf("%s(%d,%d) ", moduleName, pid, routineId))

	// req
	b.WriteString(fmt.Sprintf("<%s> ", getLogHint()))

	// 字体颜色
	colorStdout, err := utils.GetColorStdout(color)
	if err != nil {
		return "", err
	}
	b.WriteString(colorStdout)

	// 时间
	b.WriteString(time.Now().Format("2006-01-02T15:04:05"))
	b.WriteString(fmt.Sprintf("%04d", time.Now().Nanosecond()/100000))

	// 日志等级
	levelStr, err := transferLevelToStr(level)
	if err != nil {
		return "", err
	}
	b.WriteString(" ")
	b.WriteString(levelStr)
	b.WriteString(" ")

	// skip是层数，调用Caller函数外层的函数。1代表上次，2代表上上层，一般我们需要定位的也就是行数line跟file文件名
	pc, callFile, callLine, ok := runtime.Caller(s.skipCall)
	var callFuncName string
	if ok {
		// 拿到调用方法
		callFuncName = runtime.FuncForPC(pc).Name()
	}
	filePath, fileFunc := getPackageName(callFuncName)
	b.WriteString(fmt.Sprintf("%s:%d:%s ", path.Join(filePath, path.Base(callFile)), callLine, fileFunc))

	// 颜色结尾
	b.WriteString(utils.ColorEnd)
	b.WriteString(" ")

	// 文本内容
	b.WriteString(buf)
	b.WriteString("\n")

	switch s.formatType {
	case utils.FormatText:
		return b.String(), nil
	default:
		return "", errors.New("not support log format")
	}
}

func transferLevelToStr(level utils.Level) (string, error) {
	if str, ok := utils.LevelToStrMap[level]; ok {
		return str, nil
	} else {
		return "", fmt.Errorf("unknow level %d", level)
	}
}
