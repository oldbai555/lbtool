package result

import "fmt"

var _ error = (*LbErr)(nil)

type LbErr struct {
	code    uint32 `json:"code"`
	message string `json:"message"`
}

func (e *LbErr) Code() uint32 {
	return e.code
}

func (e *LbErr) Message() string {
	return e.message
}

func (e *LbErr) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.code, e.message)
}

func NewLbErr(code uint32, message string) error {
	return &LbErr{
		code:    code,
		message: message,
	}
}
