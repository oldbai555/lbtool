/**
 * @Author: zjj
 * @Date: 2025/1/13
 * @Desc:
**/

package utils

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

func GetCaller(skipCall int) string {
	pc, callFile, callLine, ok := runtime.Caller(skipCall)
	var callFuncName string
	if ok {
		// 拿到调用方法
		callFuncName = runtime.FuncForPC(pc).Name()
	}
	var getPackageName = func(f string) (filePath string, fileFunc string) {
		slashIndex := strings.LastIndex(f, "/")
		filePath = f
		if slashIndex > 0 {
			idx := strings.Index(f[slashIndex:], ".") + slashIndex
			filePath = f[:idx]
			fileFunc = f[idx+1:]
			return
		}
		return
	}
	filePath, fileFunc := getPackageName(callFuncName)
	return fmt.Sprintf("%s:%d:%s", path.Join(path.Base(filePath), path.Base(callFile)), callLine, fileFunc)
}
