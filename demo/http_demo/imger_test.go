package http_demo

import (
	"encoding/base64"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/restysdk"
	"testing"
	"time"
)

func TestCheckInterval(t *testing.T) {
	userAgent, _ := restysdk.GetRandomUserAgent()
	rsp, err := restysdk.NewRequest().SetHeader(restysdk.HeaderUserAgent, userAgent).Get("https://pps.whatsapp.net/v/t61.24694-24/315792288_609155864313005_6172241758757971000_n.jpg?stp=dst-jpg_s96x96&ccb=11-4&oh=01_AdRbFozwT7ch9sz4WPPKBqKhA8tVtr6_7ktWDLmdcIur-w&oe=63A51F35")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	rsp1, err := restysdk.NewRequest().SetHeader(restysdk.HeaderUserAgent, userAgent).Get("https://linknext-test.s3.ap-southeast-1.amazonaws.com/front_end/fe2f4b4f-fb3e-4d7b-b97b-efb7c5de9032.jpeg")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	log.Infof("%v", base64.StdEncoding.EncodeToString(rsp.Body()) == base64.StdEncoding.EncodeToString(rsp1.Body()))
}

func TestFor(t *testing.T) {
	select {
	case <-time.After(1 * time.Minute):
		log.Infof("hello world")
	}
}
