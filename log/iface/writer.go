package iface

type LogWriter interface {
	Write(p []byte) (n int, err error)
	Flush() error
}
