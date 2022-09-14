package exception

var errList []error

func Register(err ...*LbErr) {
	for _, lbErr := range err {
		errList = append(errList, lbErr)
	}
}

func GetErrCode(err error) uint32 {
	lbErr := err.(*LbErr)
	return lbErr.code
}

func GetErrMessage(err error) string {
	lbErr := err.(*LbErr)
	return lbErr.Message()
}
