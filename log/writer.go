package log

type logWriter interface {
	Write(level Level, buf string) error
	Flush() error
}
