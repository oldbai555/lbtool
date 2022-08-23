package log

var env string

var logger *Logger

func init() {
	if env == "" {
		env = DEV
	}
	logger = NewDefaultLogger(env)
}

func SetUpLogger(e string) {
	env = e
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(nil, format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(nil, format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(nil, format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(nil, format, args...)
}
