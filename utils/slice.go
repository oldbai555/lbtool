package utils

import (
	"fmt"
	"github.com/oldbai555/lbtool/extpkg/pie/pie"
	"reflect"
)

// PluckStructField2IntList 将结构体字段摘取出来转换成数组
func PluckStructField2IntList(list interface{}, fieldName string) pie.Ints {
	var result []int
	val := reflect.ValueOf(list)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			// 拿到元素
			elem := val.Index(i)
			// 指针需要进一步转换
			if elem.Kind() == reflect.Ptr {
				// Elem()获取地址指向的值
				elem = elem.Elem()
			}
			// 判断是否为空，为空就跳过
			if !elem.IsValid() {
				continue
			}
			// 判断是否为结构体
			if elem.Kind() != reflect.Struct {
				panic(any("element not struct"))
			}
			// 通过字段名称拿到这个值
			f := elem.FieldByName(fieldName)
			if !f.IsValid() {
				panic(any(fmt.Sprintf("struct missed field %s", fieldName)))
			}
			// 判断值的类型
			if f.Kind() != reflect.Int {
				panic(any(fmt.Sprintf("struct element %s type required int", fieldName)))
			}
			// 加入list中
			result = append(result, int(f.Int()))
		}
	default:
		panic(any("required list of struct type"))
	}
	return result
}

func PluckStructField2StrList(list interface{}, fieldName string) pie.Strings {
	var result []string
	val := reflect.ValueOf(list)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			// 拿到元素
			elem := val.Index(i)
			// 指针需要进一步转换
			if elem.Kind() == reflect.Ptr {
				// Elem()获取地址指向的值
				elem = elem.Elem()
			}
			// 判断是否为空，为空就跳过
			if !elem.IsValid() {
				continue
			}
			// 判断是否为结构体
			if elem.Kind() != reflect.Struct {
				panic(any("element not struct"))
			}
			// 通过字段名称拿到这个值
			f := elem.FieldByName(fieldName)
			if !f.IsValid() {
				panic(any(fmt.Sprintf("struct missed field %s", fieldName)))
			}
			// 判断值的类型
			if f.Kind() != reflect.String {
				panic(any(fmt.Sprintf("struct element %s type required int", fieldName)))
			}
			// 加入list中
			result = append(result, f.String())
		}
	default:
		panic(any("required list of struct type"))
	}
	return result
}

func UniqueSliceV2(s interface{}) interface{} {
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Slice {
		panic(any(fmt.Sprintf("s required slice, but got %v", t)))
	}

	vo := reflect.ValueOf(s)

	if vo.Len() < 2 {
		return s
	}

	res := reflect.MakeSlice(t, 0, vo.Len())
	m := map[interface{}]struct{}{}
	for i := 0; i < vo.Len(); i++ {
		el := vo.Index(i)
		eli := el.Interface()
		if _, ok := m[eli]; !ok {
			res = reflect.Append(res, el)
			m[eli] = struct{}{}
		}
	}

	return res.Interface()
}
