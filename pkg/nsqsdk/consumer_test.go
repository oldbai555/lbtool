package nsqsdk

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"testing"
	"time"
)

func TestNewConsumer(t *testing.T) {
	consumer := NewConsumer(nil, "test-channel", "127.0.0.1:4161")
	consumer.AddFunc("test", HandleMessage)
	err := consumer.Start()
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	defer consumer.Stop()
	time.Sleep(1000 * time.Millisecond)
}

// HandleMessage 处理消息
func HandleMessage(msg *nsq.Message) error {
	return Process(msg, func(data interface{}) error {
		fmt.Println("receive", msg.NSQDAddress, "message:", data)
		return nil
	})
}
