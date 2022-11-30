package webtool

import (
	"fmt"
	"github.com/oldbai555/lbtool/extpkg/lbconf/bconf"
	"github.com/oldbai555/lbtool/extpkg/lblog"
)

type Option func(tool *WebTool)

func OptionWithOrm(apollo bconf.Config, dto ...interface{}) Option {
	return func(tool *WebTool) {
		gorm := &GormMysqlConf{}
		err := gorm.InitConf(apollo)
		if err != nil {
			panic(fmt.Sprintf("err:%v", err))
		}
		err = gorm.GenConfTool(tool, dto...)
		if err != nil {
			panic(fmt.Sprintf("err:%v", err))
		}
	}
}

func OptionWithRdb(apollo bconf.Config) Option {
	return func(tool *WebTool) {
		rdb := &RedisConf{}
		err := rdb.InitConf(apollo)
		if err != nil {
			panic(fmt.Sprintf("err:%v", err))
		}
		err = rdb.GenConfTool(tool)
		if err != nil {
			panic(fmt.Sprintf("err:%v", err))
		}
	}
}

func OptionWithLog() Option {
	return func(tool *WebTool) {
		// 初始化内置日志服务
		lblog.NewLogger(lblog.SetWriteFile(true))
		tool.Log = lblog.GetLogger()
	}
}
