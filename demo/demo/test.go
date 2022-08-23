package main

import (
	"fmt"
	"time"
)

func main() {
	// 2.验证timer只能响应1次
	timer3 := time.NewTimer(2 * time.Second)

	fmt.Printf("2秒到,%v", <-timer3.C)
}
