package delayqueue

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/delayqueue/base"
	"github.com/oldbai555/lbtool/pkg/routine"
	"github.com/oldbai555/lbtool/utils"
	"math"
	"sync"
	"time"
)

var (
	_stopSignal = make(chan interface{}, 1)
	_locker     = sync.Mutex{}
	_handlers   = make(map[string]HandlerFunc)
	waitGroup   sync.WaitGroup
)

type HandlerFunc func(job *base.Job) error
type ConsumeFunc func(topic base.Topic, handler func(job *base.Job) error)

// StartConsumer 启动消费者
func StartConsumer(topics []*base.Topic) {
	log.Infof("Starting consumer")
	for _, topic := range topics {
		routine.Go(context.Background(), func(ctx context.Context) error {
			waitGroup.Add(1)
			defer waitGroup.Done()
			consume(topic)
			return nil
		})
	}
	log.Infof("Start consumer successfully")
}

// 开始消费
func consume(topic *base.Topic) {
	for {
		select {
		case <-_stopSignal:
			// 停止消费需要设计一下
			log.Infof("==========stop consume==========")
			goto exit
		default:
			// 具体消费，如果耗时大，那么可以利用管道的缓冲区进行缓冲任务执行
			job, err := listen(topic)
			if err == redis.Nil {
				continue
			}
			if err != nil {
				log.Errorf("delay_queue.Listen failed,err is %v", err)
				continue
			}
			// 没有任务，redis阻塞超时
			if job.Id == "" {
				log.Warnf("job.Id is nil,job is %v", job)
				continue
			}

			err = _handlers[topic.Name](job)
			if err != nil {
				log.Errorf("handle job failed,err is %v", err)
				job.FailedCount++

				// 任务失败时，等待时间指数级增长，最大15分钟间隔
				delay := time.Second * time.Duration(math.Pow(2, float64(job.FailedCount)))
				if delay > time.Minute*15 {
					delay = time.Minute * 15
				}

				job.ExecuteAt = time.Now().Add(delay).Unix()
				job.Topic = topic
				err = replaceJob(job)
				if err != nil {
					log.Errorf("delay_queue.Update failed,err is %v", err)
				}
				continue
			}

			err = delJob(job.Id)
			if err != nil {
				log.Errorf("delay_queue.Remove failed,err is %v", err)
				continue
			}

			// 定时任务轮询监听
			if job.Ttr > 0 {
				// 这个是定时任务的逻辑 - 得到下一次执行的时间，重新放回桶里
				newJob := base.NewJob(
					base.WithJobTTR(job.Ttr),
					base.WithJobExecuteAt(time.Now().In(utils.PRCLocation).Unix()+job.Ttr),
					base.WithJobData(job.Data),
					base.WithJobId(job.Id),
					base.WithJobTopic(job.Topic),
				)
				log.Infof("job is scheduled job,next do job is %d,job is %v", job.ExecuteAt, job)
				err = addJob2Bucket(newJob)
				if err != nil {
					log.Errorf("err:%v", err)
					continue
				}
			}

		}
	}
exit:
	log.Infof("====== close consume,topic is %v ======", topic)
	// 通知其他线程
	_stopSignal <- 1
}

// listen 轮询就绪队列获取job,正常消费消息
func listen(topic *base.Topic) (job *base.Job, err error) {
	jobId, err := blockPopFromReadyQueue(topic.Name)
	if err != nil {
		return
	}

	// 队列为空
	if jobId == "" {
		return
	}

	// 获取job元信息
	job, err = getJob(jobId)
	if err != nil {
		return
	}
	return
}

func StopConsume() {
	_stopSignal <- 1
	waitGroup.Wait()
	close(_stopSignal)
}
