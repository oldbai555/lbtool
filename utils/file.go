package utils

import (
	"log"
	"os"
)

// HasDir 判断文件夹是否存在
func HasDir(path string) (bool, error) {
	_, _err := os.Stat(path)
	if _err == nil {
		return true, nil
	}
	if os.IsNotExist(_err) {
		return false, nil
	}
	return false, _err
}

// CreateDir 创建文件夹
func CreateDir(path string) {
	log.Printf("path: %s\n", path)
	_exist, _err := HasDir(path)
	if _err != nil {
		log.Printf("获取文件夹异常 -> %v\n", _err)
		return
	}
	if _exist {
		log.Println("文件夹已存在！")
	} else {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Printf("创建目录异常 -> %v\n", err)
		} else {
			log.Println("创建成功!")
		}
	}
}
