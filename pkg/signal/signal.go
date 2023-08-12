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
var regList []func(signal os.Signal) error

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

func GetSignalChan() chan os.Signal {
	once.Do(func() {
		// Block, otherwise the listener's goroutine will exit when the main Go exits (阻塞,否则主Go退出， listenner的go将会退出)
		signalChan = make(chan os.Signal, 1)
		// 初始化监听的信号
		initReg()
	})
	return signalChan
}

func Reg(fn func(signal os.Signal) error) {
	regList = append(regList, fn)
}

func RegV2(fn func(signal os.Signal) error) {
	routine.GoV2(func() error {
		c := make(chan os.Signal, 1)
		ss := []os.Signal{
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
			syscall.SIGKILL,
		}
		signal.Notify(c, ss...)
		err := fn(<-c)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		return nil
	})
}

// Do 执行信号结束后的方法
func Do() {
	if len(regList) == 0 {
		return
	}
	routine.GoV2(func() error {
		v := <-GetSignalChan()
		for i := range regList {
			err := regList[i](v)
			if err != nil {
				log.Errorf("err:%v", err)
				continue
			}
		}
		return nil
	})
}
