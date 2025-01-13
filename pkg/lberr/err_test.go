/**
 * @Author: zjj
 * @Date: 2025/1/13
 * @Desc:
**/

package lberr

import "testing"

func TestWrapByCall(t *testing.T) {
	err := NewInvalidArg("111")
	err = WrapByCall(err)
	t.Logf("err:%v", err)
}
