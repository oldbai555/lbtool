package delay_queue

import (
	"fmt"
	"github.com/oldbai555/lb/log"
	"github.com/oldbai555/lb/pkg/exception"
)

// Get 查询Job
func Get(jobID string) (job Job, err error) {
	job, err = getJob(jobID)
	if err != nil {
		return
	}

	// 消息不存在, 可能已被删除
	if job.ID == "" {
		return
	}

	return
}

// Add 添加一个Job到队列中
func Add(job Job) error {
	if job.ID == "" || job.Topic == "" || job.ExecuteAt < 0 {
		return exception.NewInvalidArg("invalid job")
	}

	err := putJob(job.ID, job)
	if err != nil {
		log.Infof("添加job到job pool失败# putJob job is %+v,err is %s", job, err.Error())
		return err
	}

	// 拿到一个随机 bucket
	randomBucketName := <-bucketNameChan
	err = pushToBucket(randomBucketName, job.ExecuteAt, job.ID)

	if err != nil {
		log.Infof("添加job到bucket失败# pushToBucket job-%+v#%s", job, err.Error())
		return err
	}

	return nil
}

// Update 更新一个Job
func Update(job Job) (err error) {
	if job.ID == "" || job.Topic == "" || job.ExecuteAt < 0 {
		return exception.NewInvalidArg("invalid job")
	}

	err = Remove(job.ID)
	if err != nil {
		return exception.NewErr(exception.ErrDelayQueueOptErr, fmt.Sprintf("Remove job failed,err is %v", err))
	}

	err = Add(job)
	if err != nil {
		return exception.NewErr(exception.ErrDelayQueueOptErr, fmt.Sprintf("Remove job failed,err is %v", err))
	}

	return
}

// Remove 删除Job
func Remove(jobID string) error {
	return removeJob(jobID)
}
