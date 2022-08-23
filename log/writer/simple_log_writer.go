package writer

import (
	"context"
	"fmt"
	"github.com/oldbai555/comm"
	fmt2 "github.com/oldbai555/log/fmt"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

const (
	DefaultMaxFileSize = UnitMB
	UnitB              = 1
	UnitKB             = 1024 * UnitB
	UnitMB             = 1024 * UnitKB

	UnitSeconds = 1
	UnitMinutes = 60 * UnitSeconds
	UnitHour    = 60 * UnitMinutes

	DefaultChannelNumber = 1
)

func NewDefaultSimpleLoggerWriter(e string) *SimpleLoggerWriter {
	var defaultBaseDir = "./log"
	if ex, err := os.Executable(); err == nil {
		defaultBaseDir = filepath.Dir(ex) + "/log"
	}

	writer := SimpleLoggerWriter{
		env:                      e,
		baseDir:                  defaultBaseDir,
		maxFileSize:              DefaultMaxFileSize,
		checkFileFullIntervalSec: UnitSeconds * 5,
		fmt:                      fmt2.NewDefaultSimpleFormatter(),
		bufCh:                    make(chan []byte, DefaultChannelNumber),
		flushSignChan:            make(chan struct{}, DefaultChannelNumber),
		flushDoneSignChan:        make(chan error, DefaultChannelNumber),
	}
	go func() {
		err := writer.LoopDoLogic()
		panic(any(err))
	}()
	writer.isFlushing.Store(false)
	return &writer
}

// SimpleLoggerWriter 写日志
type SimpleLoggerWriter struct {
	fp                       *os.File
	env                      string
	baseDir                  string
	maxFileSize              int64
	checkFileFullIntervalSec int64  // 间隔 - 检查文件大小
	lastCheckIsFullAt        int64  // 上一次检查文件大小时间
	isFileFull               bool   // 文件是否已经满了
	currentFileName          string // 当前文件名
	fmt                      fmt2.Formatter
	openCurrentFileTime      *time.Time // 打开文件时间
	bufCh                    chan []byte
	isFlushing               atomic.Value
	flushSignChan            chan struct{} // 结束 flush 信号
	flushDoneSignChan        chan error    // 接收 flush 错误
}

// Write 写日志
func (s *SimpleLoggerWriter) Write(ctx context.Context, level fmt2.Level, format string, args ...interface{}) error {
	stdoutColor, ok := fmt2.LevelToStdoutColorMap[level]
	if !ok {
		stdoutColor = fmt2.ColorNil
	}
	logContent, err := s.fmt.Sprintf(ctx, level, stdoutColor, format, args...)
	if err != nil {
		return err
	}

	s.bufCh <- []byte(logContent)
	if s.env != common.PROD {
		log.Println(logContent)
	}
	return nil
}

// LoopDoLogic 循环执行写日志逻辑
func (s *SimpleLoggerWriter) LoopDoLogic() error {
	doWriteMoreAsPossible := func(buf []byte) error {
		for {
			var moreBuf []byte
			select {
			case moreBuf = <-s.bufCh:
				buf = append(buf, moreBuf...)
			default:
			}

			if moreBuf == nil {
				break
			}
		}

		if len(buf) == 0 {
			return nil
		}

		if err := s.tryOpenNewFile(); err != nil {
			return err
		}

		if isFull, err := s.checkFileIsFull(); err != nil {
			return err
		} else if isFull {
			fmt.Printf("log file %s is overflow max size %d bytes.", s.currentFileName, s.maxFileSize)
			return nil
		}

		bufLen := len(buf)
		var totalWrittenBytes int
		for {
			n, err := s.fp.Write(buf)
			if err != nil {
				return err
			}
			totalWrittenBytes += n
			if totalWrittenBytes >= bufLen {
				break
			}
		}

		return nil
	}

	for {
		select {
		case buf := <-s.bufCh:
			if err := doWriteMoreAsPossible(buf); err != nil {
				return err
			}
		case _ = <-s.flushSignChan:
			if err := doWriteMoreAsPossible([]byte{}); err != nil {
				s.finishFlush(err)
				break
			}
			if err := s.fp.Sync(); err != nil {
				s.finishFlush(err)
				break
			}
			s.finishFlush(nil)
		}
	}
}

// checkFileIsFull 检查文件是否满了
func (s *SimpleLoggerWriter) checkFileIsFull() (bool, error) {

	// 检查时间间隔
	if s.lastCheckIsFullAt+s.checkFileFullIntervalSec < time.Now().Unix() {
		return s.isFileFull, nil
	}

	fileInfo, err := s.fp.Stat()
	if err != nil {
		return false, err
	}

	s.isFileFull = fileInfo.Size() >= s.maxFileSize
	s.lastCheckIsFullAt = time.Now().Unix()

	return s.isFileFull, nil
}

// tryOpenNewFile 尝试开启新文件
func (s *SimpleLoggerWriter) tryOpenNewFile() error {
	var err error
	fileName := fmt.Sprintf("%s.log", time.Now().Format("2006010215"))

	if s.fp == nil {
		if _, err = os.Stat(s.baseDir); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
			if err = os.MkdirAll(s.baseDir, 0755); err != nil {
				return err
			}
		}
	}

	if s.fp, err = os.OpenFile(s.baseDir+"/"+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755); err != nil {
		return err
	}

	openFileTime := time.Now()
	s.openCurrentFileTime = &openFileTime
	s.isFileFull = false
	s.lastCheckIsFullAt = 0
	s.currentFileName = fileName

	return nil
}

// isFlushingNow 是否正在刷缓冲区
func (s *SimpleLoggerWriter) isFlushingNow() bool {
	return s.isFlushing.Load().(bool)
}

// Flush 刷缓冲区
func (s *SimpleLoggerWriter) Flush() error {
	s.isFlushing.Store(true)
	s.flushSignChan <- struct{}{}
	return <-s.flushDoneSignChan
}

// finishFlush 结束刷缓冲区
func (s *SimpleLoggerWriter) finishFlush(err error) {
	s.isFlushing.Store(false)
	s.flushDoneSignChan <- err
}

var _ LoggerWriter = (*SimpleLoggerWriter)(nil)
