package signal

import (
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/routine"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var signalChan chan os.Signal
var once sync.Once

func initReg() {
	list := []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	}
	// Listen for specified signals: ctrl+c or kill signal (监听指定信号 ctrl+c kill信号)
	signal.Notify(signalChan, list...)
}

func GetSignal() chan os.Signal {
	once.Do(func() {
		// Block, otherwise the listener's goroutine will exit when the main Go exits (阻塞,否则主Go退出， listenner的go将会退出)
		signalChan = make(chan os.Signal, 1)
		// 初始化监听的信号
		initReg()
	})
	return signalChan
}

// Do 执行信号结束后的方法
func Do(fn func(signal os.Signal) error) {
	routine.GoV2(func() error {
		v := <-GetSignal()
		err := fn(v)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		return nil
	})
}
