package main

import (
	"github.com/oldbai555/lbtool/demo/demo/config"
	"github.com/oldbai555/lbtool/log"
	"testing"
)

func Test1(t *testing.T) {
	config.M.Set("demo", "1")
	log.Infof("val is %s", config.M.Get("demo"))
	select {}
}

func Test2(t *testing.T) {
	config.M.Set("demo2", "2")
	log.Infof("val is %s", config.M.Get("demo"))
	select {}
}


