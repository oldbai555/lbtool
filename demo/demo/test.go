package main

import (
	"fmt"
	"strconv"
)

func main() {
	// 2.验证timer只能响应1次
	// timer3 := time.NewTimer(2 * time.Second)
	//
	// fmt.Printf("2秒到,%v", <-timer3.C)

	s := "23.66"
	float, _ := strconv.ParseFloat(s, 10)
	fmt.Println(float)

}

// GetLastXStr 获取最后几个字符
// prefixStr 剩下的字符
// suffixStr 最后几个字符
func GetLastXStr(str string, lastLen int) (prefixStr string, suffixStr string) {
	rs := []rune(str)
	return string(rs[:len(rs)-lastLen]), string(rs[len(rs)-lastLen:])
}
