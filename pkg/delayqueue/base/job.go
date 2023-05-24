package base

// Job 任务
type Job struct {
	// job 唯一标识
	Id string `json:"id,omitempty"`
	// 内容主体
	Data []byte `json:"data,omitempty"`

	// 所属 Topic 的 name
	Topic *Topic `json:"topic,omitempty"`
	// 失败次数
	FailedCount uint32 `json:"failed_count,omitempty"`

	// 预定执行时间 为0表示立即执行, 单位秒
	ExecuteAt int64 `json:"execute_at,omitempty"`
	// 轮询间隔 大于0表示作为定时任务, 单位秒
	Ttr int64 `json:"ttr,omitempty"`

	// 创建时间
	CreatedAt uint32 `json:"created_at,omitempty"`
}
