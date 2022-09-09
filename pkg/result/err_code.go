package result

const (
	SUCCESS             = 0
	FAILURE             = 400
	ErrInvalidArg       = 1001
	ErrOrmTableNotExist = 10001
)

var (
	Success = NewLbErr(SUCCESS, "ok")
)
