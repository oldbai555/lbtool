package comm

import (
	"reflect"
	"strings"
)

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
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
		n := relType.Field(i).Name
		if skipMap[n] {
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

// Camel2UnderScore 驼峰转下划线
func Camel2UnderScore(name string) string {
	// 检查大写字母
	var checkA2Z = func(s uint8) bool {
		return 'A' <= s && s <= 'Z'
	}

	var posList []int
	// 不从第一个字母开始
	i := 1
	for i < len(name) {
		if checkA2Z(name[i]) {
			// 记录需要转换的下标
			posList = append(posList, i)
			i++
			// 找到下一个小写字母，作为这个单词的结尾
			for i < len(name) && name[i] >= 'A' && name[i] <= 'Z' {
				i++
			}
		} else {
			i++
		}
	}
	// 全部转换为小写
	lower := strings.ToLower(name)
	if len(posList) == 0 {
		return lower
	}
	b := strings.Builder{}
	left := 0
	// 遍历需要转换的单词开头 进行转换
	for _, right := range posList {
		b.WriteString(lower[left:right])
		b.WriteByte('_')
		left = right
	}
	b.WriteString(lower[left:])
	return b.String()
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
