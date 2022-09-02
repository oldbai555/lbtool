package result

const (
	SUCCESS = 0
	FAILURE = 400
)

var (
	Success = NewLbErr(SUCCESS, "ok")
)
