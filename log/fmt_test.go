package log

import (
	"testing"
	"time"
)

func Hello() {
	Infof("hello %d, %d", time.Now().Unix(), time.Now().Unix())
}

func TestSimpleTraceFormatter_Sprintf(t *testing.T) {
	for i := 0; i < 100; i++ {
		Hello()
	}
	time.Sleep(15 * time.Second)
}
