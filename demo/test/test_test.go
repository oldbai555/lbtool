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
