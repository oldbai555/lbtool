package delay_queue

import (
	"context"
	"encoding/json"
	"github.com/oldbai555/lb/utils"
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

	// TTR 轮询间隔，非0表示定时任务
	TTR int64 `json:"ttr"`

	// ExecuteAt 预定执行时间,为0表示立即执行
	ExecuteAt uint32 `json:"execute_at"`
}

// getJob 获取Job
func getJob(key string) (job Job, err error) {
	value, err := Rdb.Get(context.TODO(), key).Result()
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
	diff := getDiffTime(job.ExecuteAt) + utils.Hours
	err := Rdb.Set(context.TODO(), key, utils.JsonEncode(job), time.Hour*time.Duration(diff)).Err()
	return err
}

// setJob 更新Job
func setJob(key string, job Job) error {
	diff := getDiffTime(job.ExecuteAt) + utils.Hours
	err := Rdb.SetNX(context.TODO(), key, utils.JsonEncode(job), time.Hour*time.Duration(diff)).Err()
	return err
}

// removeJob 删除Job
func removeJob(key string) error {
	return Rdb.Del(context.Background(), key).Err()
}

// getDiffTime 获取距离现在的时间差
func getDiffTime(executeAt uint32) uint32 {
	stampNow := utils.TimeNow()
	if executeAt < stampNow {
		return 0
	}
	return executeAt - stampNow
}
