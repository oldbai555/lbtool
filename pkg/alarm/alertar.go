package alarm

import "sync"

var (
	def     *Alerter
	defOnce sync.Once
)

func Default(svrName string) *Alerter {
	defOnce.Do(func() {
		def = InitAlerter(svrName, defaultRobotId)
	})
	return def
}
