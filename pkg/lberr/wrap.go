/**
 * @Author: zjj
 * @Date: 2025/2/25
 * @Desc:
**/

package lberr

import (
	"errors"
	"github.com/oldbai555/lbtool/utils"
)

func Join(oldErr error, errList ...error) error {
	if e, ok := oldErr.(*Error); ok {
		e.join(errList...)
		return e
	}
	var errs []error
	errs = append(errs, oldErr)
	errs = append(errs, errList...)
	return errors.Join(errs...)
}

func WrapByDesc(oldErr error, format string, args ...interface{}) error {
	wrapErr := NewErr(ErrWrapError, format, args...)
	return Join(oldErr, wrapErr)
}

func Wrap(oldErr error) error {
	wrapErr := NewErr(ErrWrapError, utils.GetCaller(2))
	return Join(oldErr, wrapErr)
}
