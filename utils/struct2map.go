package utils

import (
	"reflect"
	"strings"
	"unicode"
)

// Struct2Map 将结构体转换为map
func Struct2Map(obj interface{}) map[string]interface{} {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic("input must be a struct or a pointer to a struct")
	}

	data := make(map[string]interface{})
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if !t.Field(i).IsExported() {
			continue
		}
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

// OrmStruct2Map 转换结构体到map，并允许忽略某些字段
func OrmStruct2Map(s interface{}, skip ...string) map[string]interface{} {
	m := make(map[string]interface{})
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic("input must be a struct or a pointer to a struct")
	}

	skipMap := make(map[string]bool)
	for _, s := range skip {
		skipMap[s] = true
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		if skipMap[field.Name] {
			continue
		}

		if len(field.Name) >= 3 && field.Name[:3] == "XXX" {
			continue
		}

		if field.Name == "DeletedAt" || field.Name == "CreatedAt" || field.Name == "UpdatedAt" || field.Name == "Id" {
			continue
		}

		key := Camel2UnderScoreV2(field.Name)
		m[key] = v.Field(i).Interface()
	}

	return m
}

// OrmStruct2Map4Update 过滤掉空值
func OrmStruct2Map4Update(s interface{}, skip ...string) map[string]interface{} {
	m := OrmStruct2Map(s, skip...)
	filtered := make(map[string]interface{})
	for k, v := range m {
		if !isBlank(reflect.ValueOf(v)) {
			filtered[k] = v
		}
	}
	return filtered
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
	default:
		return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
	}
}

// Camel2UnderScoreV2 将驼峰命名转换为下划线命名
func Camel2UnderScoreV2(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}
