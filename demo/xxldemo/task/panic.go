package task

import (
	"context"
	"github.com/oldbai555/lbtool/extpkg/xxl-job-executor"
)

func Panic(cxt context.Context, param *xxl.RunReq) (msg string) {
	panic(any("test panic"))
	return
}
