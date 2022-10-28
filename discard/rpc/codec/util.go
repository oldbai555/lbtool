package codec

import "io"

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json" // not implemented
)

type NewCodecFunc func(closer io.ReadWriteCloser) Codec

var NewCodecFuncMap map[Type]NewCodecFunc

type Type string

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
