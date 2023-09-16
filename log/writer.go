package log

import (
	"fmt"
	"github.com/oldbai555/lbtool/log/_interface"
	"github.com/oldbai555/lbtool/utils"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

const (
	DefaultMaxFileSize   = utils.UnitMB
	DefaultChannelNumber = 1
)

// default linux path
var defaultBaseDir string

// SetBaseDir 没啥用 目前
func SetBaseDir(dir string) {
	defaultBaseDir = dir
}

func initDir() {
	if defaultBaseDir == "" {
		defaultBaseDir = "/tmp/lb/log"
		if runtime.GOOS == "windows" {
			defaultBaseDir = "c:/log"
		}
	}
	utils.CreateDir(defaultBaseDir)
}

func newLogWriterImpl(e string) *logWriterImpl {
	initDir()
	writer := logWriterImpl{
		env:                      e,
		baseDir:                  defaultBaseDir,
		maxFileSize:              DefaultMaxFileSize,
		checkFileFullIntervalSec: utils.Seconds * 5,
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

// logWriterImpl 写日志
type logWriterImpl struct {
	fp                       *os.File
	env                      string
	baseDir                  string
	maxFileSize              int64
	checkFileFullIntervalSec int64      // 间隔 - 检查文件大小
	lastCheckIsFullAt        int64      // 上一次检查文件大小时间
	isFileFull               bool       // 文件是否已经满了
	currentFileName          string     // 当前文件名
	openCurrentFileTime      *time.Time // 打开文件时间
	bufCh                    chan []byte
	isFlushing               atomic.Value
	flushSignChan            chan struct{} // 结束 flush 信号
	flushDoneSignChan        chan error    // 接收 flush 错误
}

// Write 写日志
func (s *logWriterImpl) Write(p []byte) (n int, err error) {

	s.bufCh <- p
	if s.env != utils.PROD {
		fmt.Printf(string(p))
	}
	return len(p), nil
}

// LoopDoLogic 循环执行写日志逻辑
func (s *logWriterImpl) LoopDoLogic() error {
	// 看看需不需要追加继续写文件
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
func (s *logWriterImpl) checkFileIsFull() (bool, error) {

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
func (s *logWriterImpl) tryOpenNewFile() error {
	var err error
	fileName := fmt.Sprintf("%s.log", time.Now().Format("2006010215"))
	if moduleName != "UNKNOWN" && moduleName != "" {
		fileName = fmt.Sprintf("%s_%s", moduleName, fileName)
	}
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
func (s *logWriterImpl) isFlushingNow() bool {
	return s.isFlushing.Load().(bool)
}

// Flush 刷缓冲区
func (s *logWriterImpl) Flush() error {
	s.isFlushing.Store(true)
	s.flushSignChan <- struct{}{}
	return <-s.flushDoneSignChan
}

// finishFlush 结束刷缓冲区
func (s *logWriterImpl) finishFlush(err error) {
	s.isFlushing.Store(false)
	s.flushDoneSignChan <- err
}

var _ _interface.LogWriter = (*logWriterImpl)(nil)
