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

//var m = protojson.MarshalOptions{}

//var um = protojson.UnmarshalOptions{}

func MarshalToString(msg proto.Message) (string, error) {
	bytes, err := protojson.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func Marshal(msg proto.Message) ([]byte, error) {
	return protojson.Marshal(msg)
}

func Unmarshal(buf []byte, msg proto.Message) error {
	return protojson.Unmarshal(buf, msg)
}
