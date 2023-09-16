package signal

import (
	"github.com/oldbai555/lbtool/log"
	"os"
	"syscall"
	"testing"
)

func TestGenSignal(t *testing.T) {
	signal := GetSignalChan()
	log.Infof("hello")
	select {
	case sig := <-signal:
		log.Infof("into")
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Infof("hello")
		case os.Interrupt, os.Kill:
			log.Infof("world")
		}
	}
}
