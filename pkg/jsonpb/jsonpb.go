/**
 * @Author: zjj
 * @Date: 2024/4/12
 * @Desc:
**/

package jsonpb

import (
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/protobuf/proto"
	"io"
)

// 简单封装 jsonpb

var m = jsonpb.Marshaler{
	EmitDefaults: true,
	OrigName:     true,
}

var um = jsonpb.Unmarshaler{
	AllowUnknownFields: true,
}

func MarshalToString(msg proto.Message) (string, error) {
	return m.MarshalToString(msg)
}

func Marshal(w io.Writer, msg proto.Message) error {
	return m.Marshal(w, msg)
}

func Unmarshal(r io.Reader, msg proto.Message) error {
	return um.Unmarshal(r, msg)
}
