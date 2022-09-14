package exception

const (
	SUCCESS             = 0
	FAILURE             = 400
	ErrInvalidArg       = 1001
	ErrOrmTableNotExist = 10001
	ErrOrmNotFound      = 10002
	ErrDelayQueueOptErr = 10003
	ErrStorageOptErr    = 10004
)

var (
	Success = NewErr(SUCCESS, "ok")
)
