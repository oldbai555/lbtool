package rpc

import (
	"fmt"
	"github.com/oldbai555/rpc/test"
	"reflect"
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	var foo test.Foo
	s := newService(&foo)
	test.Assert(len(s.method) == 1, "wrong service Method, expect 1, but got %d", len(s.method))
	mType := s.method["Sum"]
	test.Assert(mType != nil, "wrong Method, Sum shouldn't nil")
}

func TestMethodType_Call(t *testing.T) {
	var foo test.Foo
	s := newService(&foo)
	mType := s.method["Sum"]

	argv := mType.newArgv()
	replyv := mType.newReplyv()
	argv.Set(reflect.ValueOf(test.Args{Num1: 1, Num2: 3}))
	err := s.call(mType, argv, replyv)
	test.Assert(err == nil && *replyv.Interface().(*int) == 4 && mType.NumCalls() == 1, "failed to call Foo.Sum")
}

// TestSelect 测试Select 和 For 配合，退出循环
func TestSelection(t *testing.T) {
	var test = make(chan string, 1)
	for i := 0; i < 20; i++ {
		var d = i
		go func() {
		End:
			for {
				select {
				case val := <-test:
					fmt.Println(fmt.Sprintf("A%d:123456%s", d, val))
					break End
				default:
					//fmt.Println(fmt.Sprintf("A%d:wait", d))
				}
			}
		}()
	}
	for i := 0; i < 20; i++ {
		test <- "oldbai"
	}
	time.Sleep(1 * time.Second)
}
