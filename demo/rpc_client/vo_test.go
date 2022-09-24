package rpc_client

import (
	"context"
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/rpc"
	"sync"
	"testing"
	"time"
)

func TestFoo_Sum(t *testing.T) {
	client, err := rpc.Dial("tcp", fmt.Sprintf(":%d", 9999))
	if err != nil {
		log.Errorf("err is : %v", err)
		return
	}
	defer func() {
		err = client.Close()
		if err != nil {
			log.Errorf("err is : %v", err)
			return
		}
	}()

	time.Sleep(time.Second)
	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer func() {
				wg.Done()
			}()
			args := &Args{Num1: i, Num2: i * i}
			var reply int
			ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)

			if caErr := client.Call(ctx, "Foo.Sum", args, &reply); caErr != nil {
				panic(any(fmt.Sprintf("err is : %v", caErr)))
			}
			log.Infof("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}
