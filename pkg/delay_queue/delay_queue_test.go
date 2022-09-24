package delay_queue

import (
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"testing"
	"time"
)

func TestListen(t *testing.T) {
	SetupDelayQueue(&DelayQueueReq{
		Host:     "127.0.0.1",
		Port:     6379,
		Password: "",
		RedisDb:  1,
		Timeout:  10,

		BucketName: "lb-bucket",
		BucketSize: 10,

		QueueName:         "lb-ready-queue",
		QueueBlockTimeout: 10,
	})

	var topicList []Topic
	for i := 0; i < 10; i++ {
		var b = i
		topic := Topic(fmt.Sprintf("lb-topic-%d", b))
		RegisterHandler(topic, func(job Job) error {
			log.Infof("%d ===> %+v", b, job)
			return nil
		})
		topicList = append(topicList, topic)

	}
	var a = func() {
		for i := 0; i < 1000; i++ {
			utils.GenRandomStr()
			job := Job{
				Topic: topicList[i%10],
				ID:    utils.StrMd5(fmt.Sprintf("hao-%d-%d", i, utils.TimeNow())),
				Body:  fmt.Sprintf("hi %d", i),
			}
			err := Add(job)
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
		}
	}
	a()
	// ta := time.NewTicker(20 * time.Second)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ta.C:
	// 			a()
	// 		}
	// 	}
	// }()
	Start()
	time.Sleep(180 * time.Second)
}
