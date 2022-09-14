package delay_queue

import (
	"context"
	"github.com/go-redis/redis/v8"
)

// BucketItem bucket(桶)中的元素
type BucketItem struct {
	// timestamp 时间戳
	timestamp int64
	// jobID 任务ID
	jobID string
}

// 添加JobId到bucket中
func pushToBucket(randomBucketName string, timestamp uint32, jobId string) error {
	z := redis.Z{
		Score:  float64(timestamp),
		Member: jobId,
	}
	return Rdb.ZAdd(context.TODO(), randomBucketName, &z).Err()
}

// 从bucket中获取延迟时间最小的JobId
func getFromBucket(bucketName string) (*BucketItem, error) {
	value, err := Rdb.ZRangeWithScores(context.Background(), bucketName, 0, 0).Result()
	if err != nil {
		return nil, err
	}
	if value == nil || len(value) == 0 {
		return nil, nil
	}

	item := &BucketItem{}
	item.timestamp = int64(value[0].Score)
	item.jobID = (value[0].Member).(string)
	return item, nil
}

// 从bucket中删除JobId
func removeFromBucket(bucketName string, jobId string) error {
	return Rdb.ZRem(context.TODO(), bucketName, jobId).Err()
}
