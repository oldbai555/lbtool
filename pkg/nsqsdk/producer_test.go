package nsqsdk

import (
	"fmt"
	"testing"
)

func TestNewProducer(t *testing.T) {
	producer, err := NewProducer("127.0.0.1:4150")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	err = producer.np.Ping()
	if err != nil {
		return
	}

	// 生产者写入nsq,10条消息，topic = "test"
	topic := "test"
	for i := 0; i < 10; i++ {
		message := fmt.Sprintf("message:%d", i)
		if producer != nil && message != "" {
			// 不能发布空串，否则会导致error
			err = producer.Pub(topic, message) // 发布消息
			if err != nil {
				fmt.Printf("producer.Publish,err : %v", err)
			}
			fmt.Println(message)
		}
	}
	fmt.Println("producer.Publish success")
}
