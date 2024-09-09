package log

import (
	"fmt"
	"github.com/oldbai555/lbtool/env"
	"github.com/oldbai555/lbtool/log/iface"
	"github.com/oldbai555/lbtool/utils"
	"os"
	"path/filepath"
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
		defaultBaseDir = "/home/work/log"
		if runtime.GOOS == "windows" {
			defaultBaseDir = "c:/work/log"
		}
	}
	utils.CreateDir(defaultBaseDir)
}

func newLogWriterImpl() *logWriterImpl {
	initDir()
	writer := logWriterImpl{
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
	baseDir                  string        // 日志存放的目录
	maxFileSize              int64         // 文件的最大上限
	checkFileFullIntervalSec int64         // 间隔 - 检查文件大小
	lastCheckIsFullAt        int64         // 上一次检查文件大小时间
	isFileFull               bool          // 文件是否已经满了
	currentFileName          string        // 当前文件名
	openCurrentFileTime      *time.Time    // 打开文件时间
	bufCh                    chan []byte   // 缓冲区
	isFlushing               atomic.Value  // 刷盘标识
	flushSignChan            chan struct{} // 结束 flush 信号
	flushDoneSignChan        chan error    // 接收 flush 错误
}

// Write 写日志
func (s *logWriterImpl) Write(p []byte) (n int, err error) {

	s.bufCh <- p
	if !env.IsRelease() {
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

		if err := s.checkAndRotateFile(); err != nil {
			return err
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

// checkAndRotateFile 检查文件大小并决定是否需要备份和创建新文件
func (s *logWriterImpl) checkAndRotateFile() error {
	// 检查时间间隔
	if s.lastCheckIsFullAt+s.checkFileFullIntervalSec < time.Now().Unix() {
		return nil
	}

	// 获取当前文件的大小
	fileInfo, err := s.fp.Stat()
	if err != nil {
		return fmt.Errorf("无法获取文件信息: %w", err)
	}

	// 如果文件大小超过最大限制，备份并创建新文件
	if fileInfo.Size() >= s.maxFileSize {
		err := s.rotateFile()
		if err != nil {
			return fmt.Errorf("备份文件失败: %w", err)
		}
	}
	return nil
}

// rotateFile 备份当前文件并创建新文件
func (s *logWriterImpl) rotateFile() error {
	// 关闭当前文件
	err := s.fp.Close()
	if err != nil {
		return fmt.Errorf("关闭文件失败: %w", err)
	}

	// 创建备份文件名，添加时间戳或递增序号
	backupFileName := s.currentFileName + "." + time.Now().Format("20060102-150405")
	backupFilePath := filepath.Join(s.baseDir, backupFileName)

	// 将当前文件重命名为备份文件
	err = os.Rename(s.currentFileName, backupFilePath)
	if err != nil {
		return fmt.Errorf("重命名文件失败: %w", err)
	}

	// 创建新的日志文件
	newFile, err := os.Create(s.currentFileName)
	if err != nil {
		return fmt.Errorf("创建新日志文件失败: %w", err)
	}
	s.fp = newFile
	return nil
}

var _ iface.LogWriter = (*logWriterImpl)(nil)
