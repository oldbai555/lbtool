package delayqueue

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/delayqueue/base"
	"github.com/oldbai555/lbtool/utils"
	"time"
)

// ExecRedisCommand 执行redis命令, 执行完成后连接自动放回连接池
func ExecRedisCommand(command string, args ...interface{}) (interface{}, error) {
	err := conf.redisClient.Do(context.TODO(), command, args).Err()
	return nil, err
}

// 获取Job
func getJob(jobID string) (job *base.Job, err error) {
	job = &base.Job{}
	value, err := conf.redisClient.Get(context.TODO(), jobID).Result()
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(value), &job)
	if err != nil {
		return
	}

	// 消息不存在, 可能已被删除
	if job.Id == "" {
		return nil, base.ErrJobNotFound
	}
	return
}

// 添加Job
func addJob(job *base.Job) error {
	diff := utils.GetDiffTime(job.ExecuteAt) + utils.Minutes
	err := conf.redisClient.Set(context.TODO(), job.Id, utils.JsonEncode(job), time.Second*time.Duration(diff)).Err()
	return err
}

// 删除Job
func delJob(jobID string) error {
	return conf.redisClient.Del(context.Background(), jobID).Err()
}

// 推送Job入Topic中
func pubJob(topic *base.Topic, data interface{}, executeAt int64) (*base.Job, error) {
	job := base.NewJob(
		base.WithJobData(data),
		base.WithJobId(utils.Md5(data)),
		base.WithJobTopic(topic),
		base.WithJobExecuteAt(executeAt),
		//base.WithJobTTR(10), // 定时任务间隔时间
	)

	if job.ExecuteAt == 0 {
		job.ExecuteAt = time.Now().Unix()
	}
	err := addJob2Bucket(job)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	return job, nil
}

func addJob2Bucket(job *base.Job) error {
	err := addJob(job)
	if err != nil {
		log.Infof("add job fail , job is %v ,err is %v", job, err)
		return err
	}

	randomBucketName := <-randomBucketNameChan
	err = pushToBucket(randomBucketName, base.NewBucketItem(
		base.WithBucketItemExecuteAt(job.ExecuteAt),
		base.WithBucketItemData(job.Id),
	))
	if err != nil {
		log.Errorf("add job to bucket fail,job is %v,bucket is %v ,err is %v", job, randomBucketName, err)
		return err
	}
	return nil
}

// 将之前的Job给替换,用 jobId 作为Key进行替换
func replaceJob(job *base.Job) error {
	err := delJob(job.Id)
	if err != nil {
		log.Errorf("err is %v", err)
		return base.ErrRemoveJob
	}

	err = addJob2Bucket(job)
	if err != nil {
		log.Errorf("err is %v", err)
		return base.ErrPubJob
	}
	return nil
}

// 添加 桶元素 到bucket中
func pushToBucket(randomBucketName string, item *base.BucketItem) error {
	z := redis.Z{
		Score:  float64(item.ExecuteAt),
		Member: item.Data,
	}
	return conf.redisClient.ZAdd(context.TODO(), randomBucketName, &z).Err()
}

// 从bucket中获取延迟时间最小的 桶元素
func getFromBucket(bucketName string) (*base.BucketItem, error) {
	value, err := conf.redisClient.ZRangeWithScores(context.Background(), bucketName, 0, 0).Result()
	if err != nil {
		return nil, err
	}
	if value == nil || len(value) == 0 {
		return nil, nil
	}

	return base.NewBucketItem(
		base.WithBucketItemData((value[0].Member).(string)),
		base.WithBucketItemExecuteAt(int64(value[0].Score)),
	), nil
}

// 从bucket中删除 桶元素
func removeFromBucket(bucketName string, data interface{}) error {
	return conf.redisClient.ZRem(context.TODO(), bucketName, data).Err()
}

// 添加 数据 到就绪队列中
// params channel 分流数据的渠道
func pushToReadyQueue(channel string, dataMsg string) error {
	readyQueueName := conf.queue.GenReadyQueueName(channel)
	return conf.redisClient.RPush(context.TODO(), readyQueueName, dataMsg).Err()
}

// 从就绪队列中 分流数据的渠道 阻塞获取 数据
func blockPopFromReadyQueue(channel string) (string, error) {
	readyQueueName := conf.queue.GenReadyQueueName(channel)
	value, err := conf.redisClient.BLPop(context.Background(), conf.blockTimeout, readyQueueName).Result()
	if err != nil {
		return "", err
	}
	if value == nil {
		return "", nil
	}
	if len(value) == 0 {
		return "", nil
	}
	// value[0] 为弹出元素的 key
	return value[1], nil
}
