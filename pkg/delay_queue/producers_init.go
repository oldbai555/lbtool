package delay_queue

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/lb/log"
	"github.com/oldbai555/lb/utils"
	"time"
)

var (
	// 每个定时器对应一个bucket
	timers []*time.Ticker
	// bucket名称chan
	bucketNameChan <-chan string
)

type DelayQueueReq struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	RedisDb  int    `json:"redis_db"`
	Timeout  int    `json:"timeout"`

	// BucketName 存储桶名,初始化后会带上后缀 -
	BucketName string `json:"bucket_name"`
	// BucketSize 存储桶数量
	BucketSize uint32 `json:"bucket_size"`

	// QueueName ready queue在redis中的键名,初始化后会带上后缀 -
	QueueName string `json:"ready_queue_name"`
	// QueueBlockTimeout 调用 blpop 阻塞超时时间, 单位秒, 修改此项, Timeout必须做相应调整,单位为秒
	QueueBlockTimeout int `json:"ready_queue_block_timeout"`
}

var delayQueueReq *DelayQueueReq

// SetupDelayQueue 初始化延时队列
func SetupDelayQueue(req *DelayQueueReq) {
	if req == nil {
		panic(any(fmt.Sprintf("DelayQueueReq is nil")))
	}
	delayQueueReq = req
	delayQueueReq.BucketName = fmt.Sprintf("%s-", delayQueueReq.BucketName)
	delayQueueReq.QueueName = fmt.Sprintf("%s-", delayQueueReq.QueueName)

	NewRedisClient(fmt.Sprintf("%s:%d", delayQueueReq.Host, delayQueueReq.Port), delayQueueReq.Password, delayQueueReq.RedisDb, time.Duration(delayQueueReq.Timeout)*time.Second)

	initTimers(delayQueueReq.BucketSize, delayQueueReq.BucketName)

	bucketNameChan = generateBucketName(delayQueueReq.BucketSize, delayQueueReq.BucketName)
}

// initTimers 初始化定时器
func initTimers(bucketSize uint32, bucketName string) {
	timers = make([]*time.Ticker, bucketSize)

	for i := 0; i < int(bucketSize); i++ {
		timers[i] = time.NewTicker(2 * time.Second)
		newBucketName := fmt.Sprintf("%s-%d", bucketName, i+1)
		go waitTicker(timers[i], newBucketName)
	}
}

// waitTicker 启动定时器开始等待执行任务
func waitTicker(timer *time.Ticker, bucketName string) {
	for {
		select {
		case t := <-timer.C:
			tickHandler(t, bucketName)
		}
	}
}

// tickHandler 扫描bucket, 取出延迟时间小于当前时间的Job
func tickHandler(t time.Time, bucketName string) {
	t = t.In(utils.PRCLocation)
	for {
		// 从桶里拿出元素
		bucketItem, err := getFromBucket(bucketName)
		if err != nil {
			log.Errorf("扫描bucket错误#bucket-%s,err is %s", bucketName, err.Error())
			return
		}

		// 集合为空
		if bucketItem == nil {
			return
		}

		// 延迟时间未到
		if bucketItem.timestamp > t.Unix() {
			log.Warnf("%s not now,expected timestamp %d, now %d", bucketItem.jobID, bucketItem.timestamp, t.Unix())
			return
		}

		// 延迟时间小于等于当前时间, 取出Job元信息并放入ready queue
		job, err := getJob(bucketItem.jobID)
		if err != nil && err != redis.Nil {
			log.Errorf("获取Job元信息失败#jobID is %s, bucket is %s, err is %s", bucketItem.jobID, bucketName, err.Error())
			continue
		}

		// job元信息不存在, 从bucket中删除
		if err == redis.Nil || job.ID == "" {
			err = removeFromBucket(bucketName, bucketItem.jobID)
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
			continue
		}

		// 再次确认元信息中delay是否小于等于当前时间
		if job.ExecuteAt > uint32(t.Unix()) {

			// 从bucket中删除旧的jobID
			err = removeFromBucket(bucketName, bucketItem.jobID)
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}

			// 拿到一个随机 bucket ,重新计算delay时间并放入 bucket 中
			randomBucketName := <-bucketNameChan
			err = pushToBucket(randomBucketName, job.ExecuteAt, bucketItem.jobID)
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}

			continue
		}

		// 延时时间已经小于当前时间了，那么准备开始执行它
		err = pushToReadyQueue(job.Topic, bucketItem.jobID)
		if err != nil {
			log.Errorf("jobID放入ready queue失败#bucket is %s,job is %+v,err is %s", bucketName, job, err.Error())
			continue
		}

		// 从bucket中删除
		err = removeFromBucket(bucketName, bucketItem.jobID)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	}
}

// generateBucketName 轮询获取bucket名称, 使job分布到不同bucket中, 提高扫描速度
func generateBucketName(bucketSize uint32, bucketName string) <-chan string {
	c := make(chan string)
	go func() {
		i := 1
		for {
			c <- fmt.Sprintf("%s-%d", bucketName, i)
			if i >= int(bucketSize) {
				i = 1
			} else {
				i++
			}
		}
	}()

	return c
}
