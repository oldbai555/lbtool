package base

import (
	"encoding/json"
	"fmt"
	"github.com/oldbai555/lbtool/utils"
)

// ====== topic ======

type TopicOption func(topic *Topic)

func NewTopic(opts ...TopicOption) *Topic {
	var topic = &Topic{}
	for _, opt := range opts {
		opt(topic)
	}
	return topic
}

func WithTopicName(name string) TopicOption {
	return func(topic *Topic) {
		topic.Name = name
	}
}

// ====== bucket ======

type BucketOption func(bucket *Bucket)

func NewBucket(opts ...BucketOption) *Bucket {
	bucket := &Bucket{
		Name: "bucket",
		Size: 1,
	}

	for _, opt := range opts {
		opt(bucket)
	}
	return bucket
}

func WithBucketSize(size uint32) BucketOption {
	return func(bucket *Bucket) {
		bucket.Size = size
	}
}

func WithBucketName(name string) BucketOption {
	return func(bucket *Bucket) {
		bucket.Name = name
	}
}

// ====== bucket-item ======

type BucketItemOption func(item *BucketItem)

func NewBucketItem(opts ...BucketItemOption) *BucketItem {
	var item = &BucketItem{
		ExecuteAt: -1,
	}
	for _, opt := range opts {
		opt(item)
	}
	return item
}

func WithBucketItemExecuteAt(executeAt int64) BucketItemOption {
	return func(item *BucketItem) {
		item.ExecuteAt = executeAt
	}
}

func WithBucketItemData(data interface{}) BucketItemOption {
	return func(item *BucketItem) {
		item.Data = data
	}
}

// ====== queue ======

type QueueOption func(queue *Queue)

func NewQueue(opts ...QueueOption) *Queue {
	var queue = Queue{
		Name: "queue",
	}
	for _, opt := range opts {
		opt(&queue)
	}
	return &queue
}

func WithQueueName(name string) QueueOption {
	return func(queue *Queue) {
		queue.Name = name
	}
}

// ====== job ======

type JobOption func(job *Job)

// NewJob 新建一个Job
func NewJob(opts ...JobOption) *Job {
	var job = &Job{
		CreatedAt: utils.TimeNow(),
	}
	for _, opt := range opts {
		opt(job)
	}
	return job
}

func WithJobTopic(topic *Topic) JobOption {
	return func(job *Job) {
		job.Topic = topic
	}
}

func WithJobData(data interface{}) JobOption {
	return func(job *Job) {
		bytes, err := json.Marshal(data)
		if err != nil {
			panic(fmt.Sprintf("err is %v", err))
			return
		}
		job.Data = bytes
	}
}

func WithJobExecuteAt(executeAt int64) JobOption {
	return func(job *Job) {
		job.ExecuteAt = executeAt
	}
}

func WithJobTTR(ttr int64) JobOption {
	return func(job *Job) {
		job.Ttr = ttr
	}
}

func WithJobId(id string) JobOption {
	return func(job *Job) {
		job.Id = id
	}
}
