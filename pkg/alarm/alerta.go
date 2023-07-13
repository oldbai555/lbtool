package alarm

import (
	"fmt"
	"github.com/oldbai555/lbtool/pkg/warning"
	"github.com/oldbai555/lbtool/utils"
	"time"
)

type Alerter struct {
	robotId string
	svrName string
	// 不使用sync.Map 是因为就算并发导致inc的计算有误也没关系,保证至少会报警一次
	inc map[string]int64
}

func InitAlerter(svrName string, robotId string) *Alerter {
	if robotId == "" {
		robotId = defaultRobotId
	}
	alerter := &Alerter{
		svrName: svrName,
		robotId: robotId,
		inc:     make(map[string]int64),
	}
	return alerter
}

func (alerter *Alerter) Alert(title string, content string) {

	alertTitleCacheKey := alerter.getAlertTitleCacheKey(title)

	now := time.Now().Unix()
	// 会出现一直只有一个
	if sendTime, ok := alerter.inc[alertTitleCacheKey]; ok {
		if now < sendTime+10 {
			alerter.inc[alertTitleCacheKey] = now
			return
		}
	}
	alerter.inc[alertTitleCacheKey] = now
	warning.ReportToFeishu(title, fmt.Sprintf("%s: %s\n", alerter.svrName, content), alerter.robotId)
	if len(alerter.inc) > 1000 {
		warning.ReportToFeishu("Alerter", fmt.Sprintf("%s\n", "Alerter Title 超过了1000个"), alerter.robotId)
		alerter.inc = make(map[string]int64)
	}
}

func (alerter *Alerter) getAlertTitleCacheKey(title string) string {
	return fmt.Sprintf("alert_%d", utils.HashStr(title))
}

func (alerter *Alerter) AlertX(title, content string, args ...interface{}) {
	if len(args) == 0 {
		alerter.Alert(title, content)
		return
	}
	alerter.Alert(title, fmt.Sprintf(content, args...))
}
