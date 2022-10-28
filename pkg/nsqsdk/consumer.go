package nsqsdk

import (
	"fmt"
	nsq "github.com/nsqio/go-nsq"
)

type nsqDConsumer struct {
	// addrs nsqlookupd
	addrs []string
	// channel 消费队列管道
	channel string
	// entries topic 和 处理方法或函数
	entries []*entry
	// clis 消费者列表
	clis []*nsq.Consumer
	// cfg 配制
	cfg *nsq.Config
}

type entry struct {
	Topic   string
	Handler nsq.Handler
}

func NewConsumer(cfg *nsq.Config, channel string, addrs ...string) *nsqDConsumer {
	if cfg == nil {
		cfg = nsq.NewConfig()
		cfg.LookupdPollInterval = TIMEOUT
	}
	o := &nsqDConsumer{
		channel: channel,
		cfg:     cfg,
		addrs:   addrs,
	}
	return o
}

func (n *nsqDConsumer) AddFunc(topic string, fn func(*nsq.Message) error) {
	n.entries = append(n.entries, &entry{
		Topic:   topic,
		Handler: nsq.HandlerFunc(fn)})
}

func (n *nsqDConsumer) AddHandler(topic string, handler nsq.Handler) {
	n.entries = append(n.entries, &entry{
		Topic:   topic,
		Handler: handler})
}

func (n *nsqDConsumer) Start() (err error) {
	for _, e := range n.entries {
		var cli *nsq.Consumer
		cli, err = nsq.NewConsumer(e.Topic, n.channel, n.cfg)
		if err != nil {
			fmt.Printf("err:%v\n", err)
			return
		}
		cli.SetLogger(nil, 0)
		cli.AddHandler(e.Handler)

		err = cli.ConnectToNSQLookupds(n.addrs)
		if err != nil {
			return
		}
		n.clis = append(n.clis, cli)
		fmt.Printf("NewConsumer Success Topic:%s, channel:%s, addrs:%v \n", e.Topic, n.channel, n.addrs)
	}
	return
}

func (n *nsqDConsumer) Stop() {
	for _, cli := range n.clis {
		cli.Stop()
	}
	for _, cli := range n.clis {
		<-cli.StopChan
	}
}
