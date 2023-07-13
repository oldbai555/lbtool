package utils

import (
	"reflect"
)

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		if !t.Field(i).IsExported() {
			continue
		}
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func OrmStruct2Map(s interface{}, skip ...string) map[string]interface{} {
	m := make(map[string]interface{})

	elem := ValueOfData(s)

	skipMap := SliceBasis2MapValueByBool(skip).(map[string]bool)

	relType := elem.Type()

	for i := 0; i < relType.NumField(); i++ {
		if !relType.Field(i).IsExported() {
			continue
		}

		if !elem.Field(i).IsNil() {
			continue
		}

		n := relType.Field(i).Name
		if skipMap[n] {
			continue
		}
		if len(n) >= 3 &&
			n[0] == 'X' && n[1] == 'X' && n[2] == 'X' {
			continue
		}
		if n == "DeletedAt" || n == "CreatedAt" || n == "UpdatedAt" || n == "Id" {
			// skip
			continue
		}
		// 转换小驼峰
		key := Camel2UnderScore(n)
		// 把值重新写入
		m[key] = elem.Field(i).Interface()
	}

	return m
}

// OrmStruct2Map4Update 对比 OrmStruct2Map 会过滤空值
func OrmStruct2Map4Update(s interface{}, skip ...string) map[string]interface{} {
	m := OrmStruct2Map(s, skip...)

	n := map[string]interface{}{}
	for k, v := range m {
		if !isBlank(reflect.ValueOf(v)) {
			n[k] = v
		}
	}

	return n
}

// isBlank 判断值是否为空
func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}
