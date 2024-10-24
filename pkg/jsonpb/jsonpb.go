/**
 * @Author: zjj
 * @Date: 2024/4/12
 * @Desc:
**/

package jsonpb

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// 简单封装 jsonpb

var m = protojson.MarshalOptions{
	UseProtoNames: true,
}

var um = protojson.UnmarshalOptions{}

func MarshalToString(msg proto.Message) (string, error) {
	bytes, err := m.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func Marshal(msg proto.Message) ([]byte, error) {
	return m.Marshal(msg)
}

func Unmarshal(buf []byte, msg proto.Message) error {
	return um.Unmarshal(buf, msg)
}
