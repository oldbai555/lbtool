package result

import "github.com/oldbai555/lb/pkg/exception"

type LbResult struct {
	code uint32      `json:"code"`
	msg  string      `json:"msg"`
	data interface{} `json:"data"`
}

func (l *LbResult) Code() uint32 {
	return l.code
}

func (l *LbResult) Msg() string {
	return l.msg
}

func (l *LbResult) Data() interface{} {
	return l.data
}

func NewLbResult(err error, data interface{}) *LbResult {
	lbErr := err.(*exception.LbErr)
	return &LbResult{
		code: lbErr.Code(),
		msg:  lbErr.Message(),
		data: data,
	}
}

func Ok(data interface{}) *LbResult {
	return NewLbResult(exception.Success, data)
}

func Fail(msg string) *LbResult {
	return NewLbResult(exception.NewErr(exception.FAILURE, msg), "")
}
