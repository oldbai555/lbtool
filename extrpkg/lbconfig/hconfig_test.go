package hconfig

import (
	"github.com/oldbai555/lb/extrpkg/lbconfig/apollo"
	"testing"
)

type ApolloTest struct {
	Test uint32 `json:"test,omitempty" yaml:"test,omitempty"`
}

// 默认 apollo 有三个端口 8070 8080 8090 ,代码连接使用8080
func TestNewHConfig_Apollo(t *testing.T) {
	c, err := apollo.NewApolloConfig(
		apollo.WithAppid("golb"),
		apollo.WithNamespace("application.yaml"),
		apollo.WithAddr("http://127.0.0.1:8080"),
		apollo.WithCluster("DEV"),
		apollo.WithSecret("e2f16987cdc44c8fab2d0926518ba9dd"),
	)
	if err != nil {
		t.Error(err)
		return
	}

	conf, err := NewHConfig(WithDataSource(c))
	if err != nil {
		t.Error(err)
		return
	}

	// 加载配置
	if err = conf.Load(); err != nil {
		t.Error(err)
		return
	}

	//读取配置
	val, err := conf.Get("application.yaml")
	if err != nil {
		t.Error(err)
		return
	}
	//
	//var app ApolloTest
	//err = val.FormatYaml(&app)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	t.Logf("val %+v\n", val.String())

	//监听配置变化
	if err = conf.Watch(func(path string, v HVal) {
		t.Logf("path %s val %+v\n", path, v.String())
	}); err != nil {
		t.Error(err)
		return
	}
	select {}
}
