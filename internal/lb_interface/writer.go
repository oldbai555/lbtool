package lb_interface

import "github.com/oldbai555/lbtool/utils"

type LogWriter interface {
	Write(level utils.Level, buf string) error
	Flush() error
}
