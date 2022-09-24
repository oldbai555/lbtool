package delay_queue

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/lbtool/log"
	"time"
)

func getReadyQueueName(topic string) string {
	return fmt.Sprintf("%s-%s", delayQueueReq.QueueName, topic)
}

// pushToReadyQueue 添加JobId到就绪队列中
func pushToReadyQueue(topic Topic, jobId string) error {
	readyQueueName := getReadyQueueName(topic.String())
	return Rdb.RPush(context.TODO(), readyQueueName, jobId).Err()
}

// blockPopFromReadyQueue 从就绪队列中阻塞获取 JobId
func blockPopFromReadyQueue(topics []string, timeout int) (string, error) {
	var args []string
	for _, topic := range topics {
		readyQueueName := getReadyQueueName(topic)
		args = append(args, readyQueueName)
	}

	value, err := Rdb.BLPop(context.Background(), time.Duration(timeout)*time.Second, args...).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
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
