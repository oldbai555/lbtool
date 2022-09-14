package delay_queue

import (
	"fmt"
	"github.com/oldbai555/lb/log"
	"github.com/oldbai555/lb/utils"
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
			log.Infof("%d ===> %v", b, job)
			return nil
		})
		topicList = append(topicList, topic)

	}
	for i := 0; i < 100; i++ {
		utils.GenRandomStr()
		job := Job{
			Topic:     topicList[i%10],
			ID:        utils.StrMd5(fmt.Sprintf("hao-%d-%d", i, utils.TimeStampNow())),
			ExecuteAt: uint32(i+10) + uint32(time.Now().In(utils.PRCLocation).Unix()),
			Body:      fmt.Sprintf("hi %d", i),
		}
		err := Add(job)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	}
	ta := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ta.C:
				Stop()
			}
		}
	}()
	Start()
	time.Sleep(60 * time.Second)
}
