package delayqueue

import (
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/delayqueue/base"
	"testing"
	"time"
)

func TestStartDelayQueue(t *testing.T) {
	delayQueue := NewDelayQueue()
	topic1 := delayQueue.RegisterHandlerWithTopic("topic1", func(job *base.Job) error {
		log.Infof("topic1 %v", job)
		return delayQueue.Process(job, func(data interface{}) error {
			log.Infof("topic1 process %v", data)
			return nil
		})
	})

	//topic2 := delayQueue.RegisterHandlerWithTopic("topic2", func(job *base.Job) error {
	//	log.Infof("topic2 %v", job)
	//	return nil
	//})
	delayQueue.StartDelayQueue()
	defer delayQueue.StopDelayQueue()
	now := time.Now()
	for i := 0; i < 1; i++ {
		err := delayQueue.PubJob(topic1, fmt.Sprintf("hello world %d,topic 1", i), now.Add(time.Duration(i)*time.Second).Unix())
		if err != nil {
			log.Errorf("err is %v", err)
		}
		//err = delayQueue.PubJob(topic2, fmt.Sprintf("hello world %d topic2", i), now.Add(time.Duration(i)*time.Second).Unix())
		//if err != nil {
		//	log.Errorf("err is %v", err)
		//}
	}
	select {}
}
