package nsqsdk

import (
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"time"
)

type nsqDProducer struct {
	np *nsq.Producer
}

func NewProducer(addr string) (*nsqDProducer, error) {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = TIMEOUT

	p, err := nsq.NewProducer(addr, cfg)
	if err != nil {
		fmt.Printf("err:%v \n", err)
		return nil, err
	}

	np := &nsqDProducer{
		np: p,
	}
	fmt.Printf("InitProducer SUCCESS addr:%s", addr)
	return np, nil
}

func (n *nsqDProducer) Pub(topic string, c interface{}) error {
	msg, err := json.Marshal(c)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return err
	}
	b, err := EncodeMsg("", 0, msg)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return err
	}

	return n.np.Publish(topic, b)
}

func (n *nsqDProducer) DelayPub(topic string, delay time.Duration, c interface{}) error {
	msg, err := json.Marshal(c)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return err
	}
	b, err := EncodeMsg("", 0, msg)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return err
	}

	return n.np.DeferredPublish(topic, delay, b)
}

func (n *nsqDProducer) Stop() {
	n.np.Stop()
}

func (n *nsqDProducer) Ping() error {
	err := n.np.Ping()
	if err != nil {
		fmt.Printf("err:%v\n", err)
		// 关闭生产者
		n.np.Stop()
		return err
	}
	return nil
}
