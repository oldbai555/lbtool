package base

// BucketItem 桶元素
type BucketItem struct {
	ExecuteAt int64       `json:"execute_at,omitempty"` // 预期执行时间 用于redis的score排序
	Data      interface{} `json:"data,omitempty"`       // 数据
}
