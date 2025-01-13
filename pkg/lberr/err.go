package lberr

import (
	"errors"
	"fmt"
	"github.com/oldbai555/lbtool/utils"
)

var _ error = (*Error)(nil)

type Error struct {
	code    int32  `json:"code"`
	message string `json:"message"`
}

func (e *Error) Code() int32 {
	return e.code
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error[%d]:[%s]", e.code, e.message)
}

func NewErr(code int32, format string, args ...interface{}) error {
	var msg = format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	return &Error{
		code:    code,
		message: msg,
	}
}

func NewInvalidArg(format string, args ...interface{}) error {
	return NewErr(ErrInvalidArg, format, args...)
}

func NewCustomErr(format string, args ...interface{}) error {
	return NewErr(ErrCustomError, format, args...)

}

func Join(errList ...error) error {
	return errors.Join(errList...)
}

func WrapByDesc(oldErr error, format string, args ...interface{}) error {
	wrapErr := NewErr(ErrWrapError, format, args...)
	return Join(oldErr, wrapErr)
}

func Wrap(oldErr error) error {
	wrapErr := NewErr(ErrWrapError, utils.GetCaller(2))
	return Join(oldErr, wrapErr)
}
