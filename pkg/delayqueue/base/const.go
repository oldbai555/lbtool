package base

import (
	"github.com/oldbai555/lbtool/pkg/lberr"
)

const (
	ErrCodeJobNotFound = 30001
	ErrCodeRemoveJob   = 30002
	ErrCodePubJob      = 30003
)

var (
	ErrJobNotFound = lberr.NewErr(ErrCodeJobNotFound, "job not found")
	ErrRemoveJob   = lberr.NewErr(ErrCodeRemoveJob, "Remove job failed")
	ErrPubJob      = lberr.NewErr(ErrCodeRemoveJob, "Pub job failed")
)
