package delay_queue

import (
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/lb/log"
	"github.com/oldbai555/lb/utils"
	"math"
	"sync"
	"time"
)

var (
	_stopSignal = make(chan interface{}, 1)
	_locker     = sync.Mutex{}
	_handlers   = make(map[Topic]func(Job) error)
)

// Listen 轮询就绪队列获取Job
func Listen(topics ...Topic) (job Job, err error) {
	jobID, err := blockPopFromReadyQueue(topicsToStrings(topics), delayQueueReq.QueueBlockTimeout)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	// 队列为空
	if jobID == "" {
		return
	}

	// 获取job元信息
	job, err = getJob(jobID)
	if err != nil {
		return
	}

	// 消息不存在, 可能已被删除
	if job.ID == "" {
		return
	}

	// 得到下一次执行的时间，重新放回桶里
	timestamp := time.Now().In(utils.PRCLocation).Unix() + job.TTR
	err = pushToBucket(<-bucketNameChan, uint32(timestamp), job.ID)
	if err != nil {
		log.Errorf("err:%v", err)
	}
	return
}

// Consume 开始消费
func Consume(topic Topic, handler func(Job) error) {
	for {
		select {
		case <-_stopSignal:
			log.Infof("stop consume")
		default:
			job, err := Listen(topic)
			if err == redis.Nil {
				continue
			}
			if err != nil {
				log.Errorf("delay_queue.Listen failed,err is %v", err)
				continue
			}
			// 没有任务，redis阻塞超时
			if job.ID == "" {
				continue
			}

			err = handler(job)
			if err != nil {
				log.Errorf("handle job failed,err is %v", err)
				job.FailedCount++
				// 任务失败时，等待时间指数级增长，最大15分钟间隔
				delay := time.Second * time.Duration(math.Pow(2, float64(job.FailedCount)))
				if delay > time.Minute*15 {
					delay = time.Minute * 15
				}
				job.ExecuteAt = uint32(time.Now().Add(delay).Unix())
				err = Update(job)
				if err != nil {
					log.Errorf("delay_queue.Update failed,err is %v", err)
				}
				continue
			}

			err = Remove(job.ID)
			if err != nil {
				log.Errorf("delay_queue.Remove failed,err is %v", err)
			}

		}
	}
}

func Stop() {
	_stopSignal <- 1
}
