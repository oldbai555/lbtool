package json

import (
	stdjson "encoding/json"
	"io"

	"github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var j = jsoniter.Config{
	EscapeHTML:             true,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
	//	UseNumber:              true,
}.Froze()

func init() {
	extra.RegisterFuzzyDecoders()
}

// Marshal 利用json-iterator进行json编码
func Marshal(v interface{}) ([]byte, error) {
	return j.Marshal(v)
}

// MarshalIndent MarshalIndent
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return j.MarshalIndent(v, prefix, indent)
}

// Unmarshal 利用json-iterator进行json解码
func Unmarshal(data []byte, v interface{}) error {
	return j.Unmarshal(data, v)
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *jsoniter.Encoder {
	return j.NewEncoder(w)
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *jsoniter.Decoder {
	return j.NewDecoder(r)
}

type RawMessage = stdjson.RawMessage
type Number = stdjson.Number
