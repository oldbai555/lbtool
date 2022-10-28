package lberr

import "fmt"

var errMap = make(map[uint32]*LbErr)

func Register(err ...*LbErr) {
	for _, lbErr := range err {
		errMap[lbErr.code] = lbErr
	}
}

func GetErrCode(err error) uint32 {
	lbErr := err.(*LbErr)
	return lbErr.code
}

func GetErrMsg(err error) string {
	lbErr := err.(*LbErr)
	return lbErr.Message()
}

func CreateLbErr(code uint32) *LbErr {
	lbErr, ok := errMap[code]
	if ok {
		return lbErr
	}
	return &LbErr{code: code, message: fmt.Sprintf("unknown code %d", code)}
}
