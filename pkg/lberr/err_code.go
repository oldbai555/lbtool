package lberr

const (
	SUCCESS             = 0
	FAILURE             = 400
	ErrInvalidArg       = 1001
	ErrOrmTableNotExist = 10001
	ErrOrmNotFound      = 10002
	ErrDelayQueueOptErr = 10003
	ErrStorageOptErr    = 10004
	ErrNotFound         = 10005
	ErrCustomError      = 10007 // 自定义错误
	ErrRecordNotFound   = 10008
	ErrHttpError        = 10009
	ErrWrapError        = 10010 // 包装错误
)

var (
	Success        = NewErr(SUCCESS, "ok")
	RecordNotFound = NewErr(FAILURE, "record not found")
	HttpError      = NewErr(ErrHttpError, "http error")
)
