package base

import (
	"fmt"
)

type Queue struct {
	// ready queue 就绪队列的名字
	Name string `json:"name,omitempty"`
}

func (q *Queue) GenReadyQueueName(chanel string) string {
	return fmt.Sprintf("%s_%s", q.Name, chanel)
}
