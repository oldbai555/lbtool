package comm

import (
	"fmt"
	"reflect"
)

// 获取类型信息：reflect.TypeOf，是静态的
// 获取值信息：reflect.ValueOf，是动态的

// Slice2MapKeyByStructField 结构体Slice转换成Map，key field , val Struct
func Slice2MapKeyByStructField(list interface{}, fieldName string) interface{} {
	val := reflect.ValueOf(list)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 校验传入数据的类型
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
	default:
		panic(any("list required slice or array type"))
	}
	// 拿到数组的类型
	valType := val.Type()

	// 拿到数组的元素的类型
	elemType := valType.Elem()

	// elemType 用于声明 map 的 value 类型
	// elemValType 用于拿到结构体里对应字段的类型
	elemValType := elemType

	// 指针特殊处理
	for elemValType.Kind() == reflect.Ptr {
		elemValType = elemValType.Elem()
	}

	// 校验是否结构体
	if elemValType.Kind() != reflect.Struct {
		panic(any("element not struct"))
	}

	// 获取字段
	field, ok := elemValType.FieldByName(fieldName)
	if !ok {
		panic(any(fmt.Sprintf("field %s not found", fieldName)))
	}

	// 初始化存储的map
	resultMap := reflect.MakeMapWithSize(reflect.MapOf(field.Type, elemType), val.Len())

	// range slice or array set value to map
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i)
		elemStruct := elem
		for elemStruct.Kind() == reflect.Ptr {
			elemStruct = elemStruct.Elem()
		}

		// 如果是nil的，意味着key和value同时不存在，所以跳过不处理
		if !elemStruct.IsValid() {
			continue
		}

		if elemStruct.Kind() != reflect.Struct {
			panic(any("element not struct"))
		}

		resultMap.SetMapIndex(elemStruct.FieldByIndex(field.Index), elem)
	}

	return resultMap.Interface()
}

// SliceStruct2MapValueByBool 结构体Slice转换成Map，key field , val bool
func SliceStruct2MapValueByBool(list interface{}, fieldName string) interface{} {
	val := reflect.ValueOf(list)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 校验传入数据的类型
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
	default:
		panic(any("list required slice or array type"))
	}
	// 拿到数组的类型
	valType := val.Type()

	// 拿到数组的元素的类型
	elemType := valType.Elem()

	// 指针特殊处理
	for elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	// 校验是否结构体
	if elemType.Kind() != reflect.Struct {
		panic(any("element not struct"))
	}

	// 获取字段
	field, ok := elemType.FieldByName(fieldName)
	if !ok {
		panic(any(fmt.Sprintf("field %s not found", fieldName)))
	}

	// 初始化存储的map
	var bt bool
	resultMap := reflect.MakeMapWithSize(reflect.MapOf(field.Type, reflect.TypeOf(bt)), val.Len())

	// range slice or array set value to map
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i)
		elemStruct := elem
		for elemStruct.Kind() == reflect.Ptr {
			elemStruct = elemStruct.Elem()
		}

		// 如果是nil的，意味着key和value同时不存在，所以跳过不处理
		if !elemStruct.IsValid() {
			continue
		}

		resultMap.SetMapIndex(elemStruct.FieldByIndex(field.Index), reflect.ValueOf(true))
	}

	return resultMap.Interface()
}

// SliceBasis2MapValueByBool 基本数据类型Slice 转换成 map , key sliceVal val bool .
// 目前只支持整形和字符串
func SliceBasis2MapValueByBool(list interface{}) interface{} {
	val := reflect.ValueOf(list)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 拿到数组的元素的类型
	elemType := val.Type().Elem()

	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	// 初始化存储的map
	var bt bool
	var resultMap reflect.Value

	switch elemType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.String:
	default:
		panic(any("Elements of this type are not supported"))
	}

	resultMap = reflect.MakeMapWithSize(reflect.MapOf(elemType, reflect.TypeOf(bt)), val.Len())
	switch val.Kind() {
	case reflect.Array, reflect.Slice:

		for i := 0; i < val.Len(); i++ {
			// 拿到元素
			elem := val.Index(i)

			// 判断是否为空，为空就跳过
			if !elem.IsValid() {
				continue
			}

			resultMap.SetMapIndex(elem, reflect.ValueOf(true))
		}
	default:
		panic(any("required list of struct type"))
	}

	return resultMap.Interface()
}
