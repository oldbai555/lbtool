package iface

import (
	"github.com/oldbai555/lbtool/utils"
)

type Formatter interface {
	Sprintf(level utils.Level, color utils.Color, buf string) (string, error)
}
