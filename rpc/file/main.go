package main

import (
	"bufio"
	"fmt"
	"github.com/oldbai555/lb/extrpkg/pie/pie"
	"io"
	"os"
)

func main() {
	file, err := os.Open("E:\\gbai\\rpc\\file\\txt.txt")
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	createFile, err := os.Create("0628-群积分兑换活动")
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	// 写入文件
	bufio.NewWriter(createFile)
	//response, err := http.Get("hello")
	//if err != nil {
	//	fmt.Errorf("err:%v", err)
	//	return
	//}
	// 缓冲区
	//buf := new(bytes.Buffer)
	// 写对象
	//writer := bufio.NewWriter(buf)
	//reader := bufio.NewReader(response.Body)
	reader := bufio.NewReader(file)
	var i = 0
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			fmt.Printf("err:%v\n", err)
			break
		}
		if err != nil {
			fmt.Printf("err:%v\n", err)
			return
		}
		i++
		fmt.Println(fmt.Sprintf("%d - %s", i, string(line)))
		fmt.Println(pie.Ints([]int{1, 3, 4, 5}).Unique())
	}
}
