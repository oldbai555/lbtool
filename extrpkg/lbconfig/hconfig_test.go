package lbconfig

import (
	"github.com/oldbai555/lb/extrpkg/lbconfig/apollo"
	"github.com/oldbai555/lb/extrpkg/lbconfig/domain"
	"testing"
)

type ApolloTest struct {
	Test Str `json:"test,omitempty" yaml:"test,omitempty"`
}
type Str struct {
	Hello string `json:"hello,omitempty" yaml:"hello,omitempty"`
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
	val, err := conf.Get("test")
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
	t.Logf("val %+v\n", val)

	//监听配置变化
	if err = conf.Watch(func(path string, v domain.Val) {
		t.Logf("path %s val %+v\n", path, v)
		va, aerr := conf.Get("test")
		if aerr != nil {
			t.Error(aerr)
			return
		}
		t.Logf("va %+v\n", va)
	}); err != nil {
		t.Error(err)
		return
	}
	select {}
}
