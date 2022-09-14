package utils

import (
	"github.com/json-iterator/go"
	"log"
)

// JsonEncode 将任意数据JSON序列化
func JsonEncode(data interface{}) string {
	jsonStr, err := jsoniter.MarshalToString(data)
	if err != nil {
		log.Printf("JsonEncode failed ,err is %v,data is %v \n", err, data)
		return ""
	}
	return jsonStr
}
