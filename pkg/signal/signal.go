package signal

import (
	"os"
	"os/signal"
	"syscall"
)

func GenSignal() chan os.Signal {
	// Block, otherwise the listener's goroutine will exit when the main Go exits (阻塞,否则主Go退出， listenner的go将会退出)
	c := make(chan os.Signal, 1)
	// Listen for specified signals: ctrl+c or kill signal (监听指定信号 ctrl+c kill信号)
	signal.Notify(c, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	return c
}
