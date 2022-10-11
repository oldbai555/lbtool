package main

import "fmt"

const (
	Like = 1 << iota
	Collect
	Comment
)

func main() {
	ability := Like | Collect | Comment

	fmt.Printf("%b\n", ability) // 111

	fmt.Println((ability & Like) == Like)       // true
	fmt.Println((ability & Collect) == Collect) // true
	fmt.Println((ability & Comment) == Comment) // true

	fmt.Println(false && false)
}
