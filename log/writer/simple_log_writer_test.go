package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSimpleLoggerWriter_tryOpenNewFile(t *testing.T) {
	var defaultBaseDir = "./log"
	if ex, err := os.Executable(); err == nil {
		defaultBaseDir = filepath.Dir(ex) + "/log"
	}

	var file *os.File
	var err error
	if file == nil {
		if _, err = os.Stat(defaultBaseDir); err != nil {
			if !os.IsNotExist(err) {
				panic(any(err))
			}
			if err = os.MkdirAll(defaultBaseDir, 0755); err != nil {
				panic(any(err))
			}
		}
	}
	fileName := fmt.Sprintf("%s.log", time.Now().Format("20060102"))

	if file, err = os.OpenFile(defaultBaseDir+"/"+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755); err != nil {
		panic(any(err))
	}
}
