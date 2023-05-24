package main

import (
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/html"
)

func main() {
	err := html.GetHtmlInfoByUrl("https://www.yuque.com/wukong-zorrm/qdoy5p/zwre52")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	log.Errorf("err is %v", err)
	testlog()
}

func testlog() {
	log.Infof("hello")
	helloworld()
}

func helloworld() {
	log.Infof("hello world")
}
