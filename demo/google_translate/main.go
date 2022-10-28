package main

import (
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/restysdk"
)

func main() {
	c := restysdk.NewRestyClient()
	str, err := c.RestyTranslate(&restysdk.TranslateReq{
		Text: "你好呀",
		Sl:   "zh-CN",
		Tl:   "en",
	})
	if err != nil {
		log.Errorf("err:%v", err)
	}
	fmt.Println(str)

	agent, err := c.GetRandomUserAgent()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	fmt.Println(agent)
}
