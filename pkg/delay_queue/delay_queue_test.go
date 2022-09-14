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
		Timeout:  10 * time.Second,

		BucketName: "lb-bucket",
		BucketSize: 10,

		QueueName:         "lb-ready-queue",
		QueueBlockTimeout: 10,
	})

	var topicList []Topic
	for i := 0; i < 10; i++ {
		topicList = append(topicList, Topic(fmt.Sprintf("lb-%d", i)))
	}
	for i := 0; i < 100; i++ {
		utils.GenRandomStr()
		job := Job{
			Topic:     topicList[i%10],
			ID:        utils.StrMd5(fmt.Sprintf("hao-%d-%d", i, utils.TimeStampNow())),
			ExecuteAt: uint32(i+10) + uint32(time.Now().In(utils.PRCLocation).Unix()),
			TTR:       5,
			Body:      fmt.Sprintf("hi %d", i),
		}
		err := Add(job)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	}
	for {
		for i := 0; i < 10; i++ {
			receivedJob, err := Listen(topicList[i])
			if err != nil {
				log.Errorf("err:%v", err)
				continue
			}
			log.Infof("receivedJob is %v", receivedJob)
		}
	}
}
