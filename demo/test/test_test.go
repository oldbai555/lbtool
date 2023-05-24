package test

import (
	"github.com/oldbai555/lbtool/log"
	"testing"
)

func TestDo(t *testing.T) {
	CallFunc()
	log.Infof("TestDo")
}

func PrintLog() {
	log.Infof("PrintLog")
}

func CallFunc() {
	PrintLog()
	log.Infof("CallFunc")
}

func TestDo1(t *testing.T) {
	log.Infof("1 << 1 = %d", 1<<1)
	log.Infof("1 << 2 = %d", 1<<2)
	log.Infof("1 << 3 = %d", 1<<3)
	log.Infof("1 | 1 << 3 = %d", 1|1<<3)
	log.Infof("2 | 1 << 3 = %d", 2|1<<3)
	log.Infof("4 | 1 << 3 = %d", 4|1<<3)
	log.Infof("8 | 1 << 3 = %d", 8|1<<3)
}
