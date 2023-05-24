package main

import (
	"bytes"
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
	"io/ioutil"
	"os"
)

func main() {
	//1.读取标准输入，接收proto 解析的文件内容，并解析成结构体
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("err is : %v\n", err)
		panic(any(err))
	}

	var req pluginpb.CodeGeneratorRequest
	err = proto.Unmarshal(input, &req)
	if err != nil {
		fmt.Printf("err is : %v\n", err)
		panic(any(err))
	}

	//2.生成插件
	opts := protogen.Options{}
	plugin, err := opts.New(&req)
	if err != nil {
		panic(any(err))
	}

	// 3.在插件 plugin.Files 就是 demo.proto 的内容了,是一个切片，每个切片元素代表一个文件内容
	// 我们只需要遍历这个文件就能获取到文件的信息了
	for _, file := range plugin.Files {
		//创建一个buf 写入生成的文件内容
		var buf bytes.Buffer

		// 写入go 文件的package名
		pkg := fmt.Sprintf("package %s", file.GoPackageName)
		buf.Write([]byte(pkg))

		//遍历消息,这个内容就是protobuf的每个消息
		for _, msg := range file.Messages {
			//接下来为每个消息生成hello 方法

			buf.Write([]byte(fmt.Sprintf(`
             func (m *%s)Hello(){

             }
             `, msg.GoIdent.GoName)))
		}

		//指定输入文件名,输出文件名为demo.foo.go
		filename := file.GeneratedFilenamePrefix + ".foo.go"
		f := plugin.NewGeneratedFile(filename, ".")

		// 将内容写入插件文件内容
		_, err = f.Write(buf.Bytes())
		if err != nil {
			fmt.Printf("err is : %v\n", err)
			panic(any(err))
		}

	}

	// 生成响应
	stdout := plugin.Response()
	out, err := proto.Marshal(stdout)
	if err != nil {
		panic(any(err))
	}

	// 将响应写回 标准输入, protoc会读取这个内容
	fmt.Fprintf(os.Stdout, string(out))
}
