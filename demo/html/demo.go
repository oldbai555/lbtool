package main

import (
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/html"
)

func main() {
	err := html.GetHtmlInfoByUrl("https://zhuanlan.zhihu.com/p/387840381")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
