package main

import (
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"github.com/panjf2000/ants/v2"
	"time"
)

func main() {
	// 第一种用法
	pool, err := ants.NewPool(10)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	defer func() {
		rErr := pool.ReleaseTimeout(10 * time.Second)
		if rErr != nil {
			log.Errorf("err:%v", rErr)
			return
		}
	}()

	for i := 0; i < 100; i++ {
		var a = i
		err = pool.Submit(func() {
			log.Infof("%d", a)
			time.Sleep(2 * time.Second)
		})
	}

	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	// 第二种用法
	poolWF, err := ants.NewPoolWithFunc(10, func(i interface{}) {
		log.Warnf("%v", i)
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	for i := 0; i < 100; i++ {
		var a = i
		err = poolWF.Invoke(fmt.Sprintf("hello + %d", a))
	}
}
