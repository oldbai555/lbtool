package fmt

import (
	"github.com/oldbai555/log"
	"testing"
	"time"
)

func Hello() {
	log.Infof("hello %d, %d", time.Now().Unix(), time.Now().Unix())
}

func TestSimpleTraceFormatter_Sprintf(t *testing.T) {
	for i := 0; i < 100; i++ {
		Hello()
	}
	time.Sleep(15 * time.Second)
}
