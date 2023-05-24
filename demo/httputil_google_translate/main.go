package main

import (
	"bytes"
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/restysdk"
	"io/ioutil"
)

func Get() {
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

	rsp, err := restysdk.NewRequest().Get("https://pps.whatsapp.net/v/t61.24694-24/309079747_143263165076512_1566853013342809496_n.jpg?ccb=11-4&oh=01_AdQ9QYPuwGhh5IhckWXHcmZGWyOQKWu3zt9bRgsZVWTqCA&oe=639FAAF6")
	// var buf bytes.Buffer
	newBuffer := bytes.NewBuffer(rsp.Body())
	ioutil.WriteFile("309079747_143263165076512_1566853013342809496_n.jpg", newBuffer.Bytes(), 0666)
}
