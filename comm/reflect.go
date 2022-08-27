package comm

import (
	"reflect"
)

// KindOfData 获取数据的类型
func KindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

// ValueOfData 获取数据,优化下指针
func ValueOfData(data interface{}) reflect.Value {
	value := reflect.ValueOf(data)
	valueType := value

	if valueType.Kind() == reflect.Ptr {
		valueType = value.Elem()
	}
	return valueType
}

// ToSlice 将slice 转换成 []internal{}
func ToSlice(arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		panic(any("toslice arr not slice"))
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret
}

// ReflectCompareFieldDesc 按传入fieldName 比较.降序
func ReflectCompareFieldDesc(i, j interface{}, fieldName string) bool { //
	valI := ValueOfData(i).FieldByName(fieldName).Interface()
	valJ := ValueOfData(j).FieldByName(fieldName).Interface()
	switch s := valI.(type) {
	case string:
		return s < valJ.(string)
	case float64:
		return s < valJ.(float64)
	case int:
		return s < valJ.(int)
	case uint32:
		return s < valJ.(uint32)
	default:
		panic(any("The type is unknown"))
	}
	return true
}

func ReflectCompareFieldAsc(i, j interface{}, fieldName string) bool { //
	valI := ValueOfData(i).FieldByName(fieldName).Interface()
	valJ := ValueOfData(j).FieldByName(fieldName).Interface()
	switch s := valI.(type) {
	case string:
		return s > valJ.(string)
	case float64:
		return s > valJ.(float64)
	case float32:
		return s > valJ.(float32)
	case int:
		return s > valJ.(int)
	case uint32:
		return s > valJ.(uint32)
	default:
		panic(any("The type is unknown"))
	}
	return false
}
