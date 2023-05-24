package delayqueue

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/delayqueue/base"
	"github.com/oldbai555/lbtool/pkg/routine"
	"time"
)

var tickers []*time.Ticker

func StartProducer(bucket *base.Bucket) {
	log.Infof("Starting producer")
	randomBucketNameChan = GenerateRandomBucketCacheKeyChan(bucket)
	tickers = make([]*time.Ticker, bucket.Size)
	for i := uint32(0); i < bucket.Size; i++ {
		bucketName := bucket.GenBucketName(i)
		ticker := time.NewTicker(time.Second)
		routine.Go(context.Background(), func(ctx context.Context) error {
			log.Infof("loading producer......")
			waitTicker(ticker, bucketName)
			return nil
		})
		tickers[i] = ticker
	}
	log.Infof("Start producer successfully")
}

// waitTicker 启动定时器开始等待执行任务
func waitTicker(ticker *time.Ticker, bucketName string) {
	for {
		select {
		case t := <-ticker.C:
			tickHandler(t, bucketName)
		}
	}
}

// tickHandler 扫描bucket, 取出延迟时间小于当前时间的Job
func tickHandler(t time.Time, bucketName string) {
	// 放循环是为了拿同一时间的Job
	for {
		// 从桶里拿出元素
		bucketItem, err := getFromBucket(bucketName)
		if err != nil {
			log.Errorf("扫描 bucket 错误 bucket is %s,err is %s", bucketName, err.Error())
			return
		}

		// 集合为空
		if bucketItem == nil {
			log.Debugf("bucket item is nil")
			return
		}

		// 延迟时间未到
		if bucketItem.ExecuteAt > t.Unix() {
			log.Warnf("%s not now,expected executeAt %d, now %d", bucketItem.Data, bucketItem.ExecuteAt, t.Unix())
			return
		}

		// 延迟时间小于等于当前时间, 取出Job元信息并放入ready queue
		job, err := getJob(fmt.Sprintf("%s", bucketItem.Data))
		if err != nil && err != redis.Nil {
			log.Errorf("获取Job元信息失败#jobID is %s, bucket is %s, err is %s", bucketItem.Data, bucketName, err.Error())
			return
		}

		// job元信息不存在, 从bucket中删除
		if err == redis.Nil || job.Id == "" {
			err = removeFromBucket(bucketName, bucketItem.Data)
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
			return
		}

		// 再次确认元信息中delay是否小于等于当前时间
		if job.ExecuteAt > t.Unix() {

			// 从bucket中删除旧的jobID
			err = removeFromBucket(bucketName, bucketItem.Data)
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}

			// 拿到一个随机 bucket ,重新计算delay时间并放入 bucket 中
			randomBucketName := <-randomBucketNameChan
			err = pushToBucket(randomBucketName, base.NewBucketItem(
				base.WithBucketItemExecuteAt(job.ExecuteAt),
				base.WithBucketItemData(bucketItem.Data),
			))
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
		}

		// 延时时间已经小于当前时间了，那么准备开始执行它
		err = pushToReadyQueue(job.Topic.Name, fmt.Sprintf("%s", bucketItem.Data))
		if err != nil {
			log.Errorf("jobID放入ready queue失败#bucket is %s,job is %+v,err is %s", bucketName, job, err.Error())
			return
		}

		// 从bucket中删除
		err = removeFromBucket(bucketName, fmt.Sprintf("%s", bucketItem.Data))
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	}

}

// GenerateRandomBucketCacheKeyChan 轮询获取bucket名称, 使job分布到不同bucket中, 提高扫描速度
func GenerateRandomBucketCacheKeyChan(bucket *base.Bucket) <-chan string {
	log.Infof("Starting random bucket")
	c := make(chan string)
	routine.Go(context.Background(), func(ctx context.Context) error {
		log.Infof("loading random bucket......")
		index := uint32(0)
		for {
			// 重置一下 避免溢出
			if index >= bucket.Size {
				index = 0
			}
			// 拿到随机桶的缓存 key
			c <- bucket.GenBucketName(index)
			index++
		}
	})
	log.Infof("Start random bucket SuccessFully")
	return c
}
