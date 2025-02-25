package lberr

import (
	"fmt"
)

var _ error = (*Error)(nil)

type Error struct {
	code    int32
	message string
	errs    []error
}

func (e *Error) Code() int32 {
	return e.code
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) Cause() error {
	return &Error{
		code:    e.code,
		message: e.message,
	}
}

func (e *Error) Error() string {
	var b []byte
	for i, err := range e.errs {
		if i > 0 {
			b = append(b, '\t')
		}
		b = append(b, err.Error()...)
	}
	appendErrorStr := string(b)
	if len(appendErrorStr) == 0 {
		if e.code == ErrWrapError {
			return e.message
		}
		return fmt.Sprintf("code: %d, msg: %s", e.code, e.message)
	}
	return fmt.Sprintf("code: %d,msg: %s\t%s", e.code, e.message, appendErrorStr)
}

func (e *Error) join(errs ...error) {
	n := 0
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return
	}
	for _, err := range errs {
		if err != nil {
			e.errs = append(e.errs, err)
		}
	}
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
