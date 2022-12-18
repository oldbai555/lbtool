package delayqueue

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/lbtool/log"
	base2 "github.com/oldbai555/lbtool/pkg/delayqueue/base"
	"time"
)

var (
	conf                 *Conf
	defaultTimeout       = time.Duration(10) * time.Second
	randomBucketNameChan <-chan string // 随机桶管道
)

type Conf struct {
	redisClient *redis.Client
	// 调用 blPop 阻塞超时时间, 单位秒, 默认十秒,修改此项, redis Timeout 必须做相应调整,单位为秒
	blockTimeout time.Duration

	bucket *base2.Bucket
	queue  *base2.Queue

	topics []*base2.Topic
}

type Option func(conf *Conf)

// NewDelayQueue 初始化延时队列配置
func NewDelayQueue(optList ...Option) *Conf {
	conf = &Conf{
		bucket:       base2.NewBucket(),
		queue:        base2.NewQueue(),
		blockTimeout: defaultTimeout,
	}
	for _, opt := range optList {
		opt(conf)
	}
	if conf.redisClient == nil {
		conf.redisClient = redis.NewClient(&redis.Options{
			Addr:        fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
			Password:    "", // no password set
			DB:          1,  // use default DB
			ReadTimeout: defaultTimeout,
		})
	}
	return conf
}

// WithRedisClient 初始化缓存客户端
func WithRedisClient(redidOpt *redis.Options) Option {
	return func(conf *Conf) {
		conf.redisClient = redis.NewClient(redidOpt)
	}
}

func WithBucket(bucket *base2.Bucket) Option {
	return func(conf *Conf) {
		conf.bucket = bucket
	}
}

func WithConf(queue *base2.Queue) Option {
	return func(conf *Conf) {
		conf.queue = queue
	}
}

func (conf *Conf) RegisterHandlerWithTopic(topicName string, handler HandlerFunc) *base2.Topic {
	log.Infof("registry handler to topic,topicName: %s", topicName)
	_locker.Lock()
	defer _locker.Unlock()
	_handlers[topicName] = handler
	topic := base2.NewTopic(
		base2.WithTopicName(topicName),
	)
	conf.topics = append(conf.topics, topic)
	return topic
}

func (conf *Conf) StartDelayQueue() {
	StartProducer(conf.bucket)
	StartConsumer(conf.topics)
}

func (conf *Conf) StopDelayQueue() {
	StopConsume()
	for _, timer := range tickers {
		timer.Stop()
	}
}

type doTopicHandlerFunc func(data interface{}) error

func (conf *Conf) Process(job *base2.Job, f doTopicHandlerFunc) error {
	var data interface{}
	err := json.Unmarshal(job.Data, &data)
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	return f(data)
}

func (conf *Conf) PubJob(topic *base2.Topic, data interface{}, executeAt int64) error {
	job, err := pubJob(topic, data, executeAt)
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	log.Infof("pub job successfully,job id is %s", job.Id)
	return nil
}
