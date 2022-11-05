package webtool

import (
	"github.com/oldbai555/lbtool/extpkg/lbconf/bconf"
)

type Option interface {
	InitConf(apollo bconf.Config) error
	GenConfTool(tool *WebTool) error
}
