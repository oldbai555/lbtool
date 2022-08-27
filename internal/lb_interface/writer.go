package lb_interface

import "github.com/oldbai555/lb/utils"

type LogWriter interface {
	Write(level utils.Level, buf string) error
	Flush() error
}
