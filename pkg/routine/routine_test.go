package routine

import (
	"testing"
	"time"
)

func TestGo(t *testing.T) {
	Go(func() error {
		panic("abc")
	})
	time.Sleep(5 * time.Second)
}
