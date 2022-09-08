package result

const (
	SUCCESS             = 0
	FAILURE             = 400
	ErrOrmTableNotExist = 10001
)

var (
	Success = NewLbErr(SUCCESS, "ok")
)
