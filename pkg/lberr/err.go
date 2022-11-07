package lberr

import "fmt"

var _ error = (*LbErr)(nil)

type LbErr struct {
	code    int32  `json:"code"`
	message string `json:"message"`
}

func (e *LbErr) Code() int32 {
	return e.code
}

func (e *LbErr) Message() string {
	return e.message
}

func (e *LbErr) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.code, e.message)
}

func NewErr(code int32, message string) error {
	return &LbErr{
		code:    code,
		message: message,
	}
}

func NewInvalidArg(format string, args ...interface{}) error {
	return &LbErr{
		code:    ErrInvalidArg,
		message: fmt.Sprintf(format, args),
	}
}

func NewCustomErr(format string, args ...interface{}) error {
	return &LbErr{
		code:    ErrCustomError,
		message: fmt.Sprintf(format, args),
	}
}
