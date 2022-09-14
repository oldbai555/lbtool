package delay_queue

import (
	"context"
	"fmt"
	"github.com/oldbai555/lb/log"
	"time"
)

// pushToReadyQueue 添加JobId到就绪队列中
func pushToReadyQueue(topic Topic, jobId string) error {
	readyQueueName := fmt.Sprintf(delayQueueReq.QueueName, topic)
	return Rdb.RPush(context.TODO(), readyQueueName, jobId).Err()
}

// blockPopFromReadyQueue 从就绪队列中阻塞获取 JobId
func blockPopFromReadyQueue(topics []string, timeout int) (string, error) {
	var args []string
	for _, topic := range topics {
		readyQueueName := fmt.Sprintf(delayQueueReq.QueueName, topic)
		args = append(args, readyQueueName)
	}

	value, err := Rdb.BLPop(context.Background(), time.Duration(timeout)*time.Second, args...).Result()
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}
	if value == nil {
		return "", nil
	}
	if len(value) == 0 {
		return "", nil
	}
	element := value[1]

	return element, nil
}
