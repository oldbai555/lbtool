package lb_interface

import (
	"github.com/oldbai555/lb/utils"
)

type Formatter interface {
	Sprintf(level utils.Level, color utils.Color, buf string) (string, error)
}
