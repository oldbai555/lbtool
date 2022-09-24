package delay_queue

import (
	"context"
	"encoding/json"
	"github.com/oldbai555/lbtool/utils"
	"time"
)

type JobPrefix string

func (o JobPrefix) String() string {
	return string(o)
}

// Job 任务
type Job struct {
	// Topic 主题
	Topic Topic `json:"topic"`
	// ID job唯一标识ID
	ID string `json:"id"`

	// FailedCount 失败次数
	FailedCount int64 `json:"failed_count"`
	// Body 内容主体
	Body string `json:"body"`

	// TTR 轮询间隔，大于0表示定时任务,单位秒
	TTR int64 `json:"ttr"`

	// ExecuteAt 预定执行时间,为0表示立即执行,时间戳
	ExecuteAt uint32 `json:"execute_at"`
}

// getJob 获取Job
func getJob(jobID string) (job Job, err error) {
	value, err := Rdb.Get(context.TODO(), jobID).Result()
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(value), &job)
	if err != nil {
		return
	}

	return
}

// putJob 添加Job
func putJob(key string, job Job) error {
	diff := getDiffTime(job.ExecuteAt) + utils.Minutes
	err := Rdb.Set(context.TODO(), key, utils.JsonEncode(job), time.Second*time.Duration(diff)).Err()
	return err
}

// removeJob 删除Job
func removeJob(jobID string) error {
	return Rdb.Del(context.Background(), jobID).Err()
}

// getDiffTime 获取距离现在的时间差,单位秒
func getDiffTime(executeAt uint32) uint32 {
	stampNow := utils.TimeNow()
	if executeAt < stampNow {
		return 0
	}
	return executeAt - stampNow
}
