package base

import (
	"fmt"
)

// Bucket 存储桶
type Bucket struct {
	Name string `json:"name,omitempty"`
	Size uint32 `json:"size,omitempty"`
}

func (b *Bucket) GenBucketName(index uint32) string {
	return fmt.Sprintf("%s_%d", b.Name, index)
}
