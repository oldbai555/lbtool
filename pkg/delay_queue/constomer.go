package delay_queue

import (
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/lb/log"
	"github.com/oldbai555/lb/utils"
	"math"
	"runtime"
	"sync"
	"time"
)

var (
	_stopSignal = make(chan interface{}, 1)
	_locker     = sync.Mutex{}
	_handlers   = make(map[Topic]HandlerFunc)
)

type HandlerFunc func(Job) error
type ConsumeFunc func(topic Topic, handler func(Job) error)

// RegisterHandler 注册方法
func RegisterHandler(topic Topic, handler HandlerFunc) {
	_locker.Lock()
	defer _locker.Unlock()
	_handlers[topic] = handler
}

// Start 启动消费者
func Start() {
	log.Infof("==========start consume==========")
	for topic, handler := range _handlers {
		go protectedRun(topic, handler, consume)
	}
}

// protectedRun 守护运行，避免宕机
func protectedRun(topic Topic, handler HandlerFunc, fn ConsumeFunc) {
	// 延迟处理的函数
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		err := recover()
		switch any(err).(type) {
		case runtime.Error: // 运行时错误
			log.Errorf("runtime error:", err)
		default: // 非运行时错误
			log.Errorf("error:", err)
		}
	}()
	fn(topic, handler)
}

// listen 轮询就绪队列获取job,正常消费消息
func listen(topics ...Topic) (job Job, err error) {
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

	// 定时任务轮询监听
	if job.TTR > 0 {
		// 这个是定时任务的逻辑 - 得到下一次执行的时间，重新放回桶里
		timestamp := time.Now().In(utils.PRCLocation).Unix() + job.TTR
		err = pushToBucket(<-bucketNameChan, uint32(timestamp), job.ID)
		if err != nil {
			log.Errorf("err:%v", err)
		}
	}

	return
}

// consume 开始消费
func consume(topic Topic, handler func(Job) error) {
	for {
		select {
		case <-_stopSignal:
			// todo 停止消费需要设计一下
			log.Infof("==========stop consume==========")
		default:
			job, err := listen(topic)
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

// Stop 停止消费
func Stop() {
	_stopSignal <- 1
}
