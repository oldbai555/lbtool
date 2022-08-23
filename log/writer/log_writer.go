package writer

import (
	"context"
	"lb/log/fmt"
)

type LoggerWriter interface {
	Write(ctx context.Context, level fmt.Level, format string, args ...interface{}) error
	Flush() error
}
