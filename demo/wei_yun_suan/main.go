package main

import "fmt"

const (
	Like = 1 << iota
	Collect
	Comment
)

// https://juejin.cn/post/7090884238588772383
func main() {
	ability := Like | Collect | Comment

	fmt.Printf("%b\n", ability) // 111

	fmt.Println((ability & Like) == Like)       // true
	fmt.Println((ability & Collect) == Collect) // true
	fmt.Println((ability & Comment) == Comment) // true

	fmt.Printf("%v\n", ability^Like)
	fmt.Printf("%v\n", ability^Like^Like)
}
