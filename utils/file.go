package utils

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
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
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Printf("创建目录异常 -> %v\n", err)
		} else {
			log.Println("创建成功!")
		}
	}
}

// GetCurrentAbPath 最终方案-获取当前程序的绝对路径
func GetCurrentAbPath() string {
	dir := getCurrentAbPathByExecutable()
	if strings.Contains(dir, getTmpDir()) {
		return getCurrentAbPathByCaller()
	}
	log.Printf("get_current_ab_path is %s \n", dir)
	return dir
}

// 获取系统临时目录，兼容go run
func getTmpDir() string {
	dir := os.Getenv("TEMP")
	if dir == "" {
		dir = os.Getenv("TMP")
	}
	res, _ := filepath.EvalSymlinks(dir)
	return res
}

// 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(2)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

// FileExists 文件是否存在
func FileExists(name string) bool {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func GetCurDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1) + "/"
}
