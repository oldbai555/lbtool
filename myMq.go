package impl

import (
	"fmt"
	"git.pinquest.cn/base/log"
	"git.pinquest.cn/qlb/brick/smq"
	"git.pinquest.cn/xingyu/hw/internal/mq"
	"git.pinquest.cn/xingyu/hw/service/hwmallconver"
)

func InitMq() (err error) {

	// 推送且自己内部消费
	_, err = smq.AddTopic(nil, &smq.AddTopicReq{
		Topic: &smq.Topic{
			Name: mq.HwShopifyOrderInfoMq.TopicName,
			SubConfig: &smq.SubConfig{
				ServiceName:     hwmallconver.ServiceName,
				ServicePath:     hwmallconver.HandlerHwShopifyOrderInfoMqCMDPath,
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
	_, err = smq.AddTopic(nil, &smq.AddTopicReq{
		Topic: &smq.Topic{
			Name: mq.HwShopifyOrderInfoMq.TopicName,
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
	_, err = smq.AddChannel(nil, &smq.AddChannelReq{
		TopicName: mq.HwShopifyCustomerInfoMq.TopicName,
		Channel: &smq.Channel{
			Name: fmt.Sprintf("%s_%s", hwmallconver.SvrName, "customer_info_mq"),
			SubConfig: &smq.SubConfig{
				ServiceName:     hwmallconver.ServiceName,
				ServicePath:     hwmallconver.HandlerHwShopifyOrderInfoMqCMDPath,
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
