package main

import (
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/template"
)

func main() {

	var funcT = &template.Function{
		Package:     "impl",
		ModelName:   "hwmallconver",
		Variable:    "HwShopifyCartInfoMq",
		Template:    initMqTemp,
		Description: "购物车信息 - 对业务",
	}

	file, err := template.GetOsFile("D:\\lb\\demo\\temp", "myMq")
	if err != nil {
		log.Errorf("err:%v", err)
	}
	err = template.GenTemplate(file, funcT)
	if err != nil {
		log.Errorf("err:%v", err)
	}
}

var initMqTemp = `
package {{.Package}}

import (
	"fmt"
	"git.pinquest.cn/base/log"
	"git.pinquest.cn/qlb/brick/smq"
	"git.pinquest.cn/xingyu/hw/internal/mq"
	"git.pinquest.cn/xingyu/hw/service/{{.ModelName}}"
)

func InitMq() (err error) {

	// 推送且自己内部消费
	// {{.Description}}
	_, err = smq.AddTopic(nil, &smq.AddTopicReq{
		Topic: &smq.Topic{
			Name: mq.{{.Variable}}.TopicName,
			SubConfig: &smq.SubConfig{
				ServiceName:     {{.ModelName}}.ServiceName,
				ServicePath:     {{.ModelName}}.Handler{{.Variable}}CMDPath,
				ConcurrentCount: 1,
				MaxRetryCount:   5,
			},
		},
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// 只推送
	// {{.Description}}
	_, err = smq.AddTopic(nil, &smq.AddTopicReq{
		Topic: &smq.Topic{
			Name: mq.{{.Variable}}.TopicName,
			SubConfig: &smq.SubConfig{
				ConcurrentCount: 10,
				MaxRetryCount:   5,
			},
			SkipConsumeTopic: true,
		},
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// 只消费
	// {{.Description}}
	_, err = smq.AddChannel(nil, &smq.AddChannelReq{
		TopicName: mq.HwShopifyCustomerInfoMq.TopicName,
		Channel: &smq.Channel{
			Name: fmt.Sprintf("%s_%s", hwmallconver.SvrName, "customer_info_mq"),
			SubConfig: &smq.SubConfig{
				ServiceName:     {{.ModelName}}.ServiceName,
				ServicePath:     {{.ModelName}}.Handler{{.Variable}}CMDPath,
				ConcurrentCount: 10,
				MaxRetryCount:   5,
			},
		},
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

}
`
