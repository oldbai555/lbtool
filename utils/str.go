package utils

import (
	"hash/fnv"
	"strings"
	"unicode"
)

// UpperFirst 首字母大写
func UpperFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// LowerFirst 首字母小写
func LowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// SubStr 截取字符
func SubStr(str string, start, end uint32) string {
	runeTitle := []rune(str)
	return string(runeTitle[start:end])
}

// UnderScore2Camel 下划线转驼峰
func UnderScore2Camel(name string) string {
	var buf []byte
	toggleUpper := true
	for i := 0; i < len(name); i++ {
		if name[i] == '_' {
			toggleUpper = true
		} else {
			c := name[i]
			if toggleUpper {
				toggleUpper = false
				if c >= 'a' && c <= 'z' {
					c = c - 'a' + 'A'
				}
			}
			if c >= '0' && c <= '9' {
				toggleUpper = true
			}
			buf = append(buf, c)
		}
	}
	return string(buf)
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

// HashStr 哈希字符串
func HashStr(s string) uint32 {
	f := fnv.New32a()
	_, _ = f.Write([]byte(s))
	return f.Sum32()
}
