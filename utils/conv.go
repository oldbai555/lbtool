package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// vec = vector 动态数组

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err == nil {
		return i
	}
	return 0
}

func Itoa(i uint64) string {
	return strconv.FormatUint(uint64(i), 10)
}
func I64toa(i int64) string {
	return strconv.FormatInt(i, 10)
}
func I32toa(i uint32) string {
	return strconv.FormatUint(uint64(i), 10)
}

func Int2Str(i interface{}) string {
	return fmt.Sprintf("%v", i)
}

func AtoUint64(s string) uint64 {
	return AToU64(s)
}

func AtoUint32(s string) uint32 {
	return uint32(AToU64(s))
}

func AtoInt64(s string) int64 {
	return int64(AToF64(s))
}

func AtoInt(s string) int {
	return int(AToF64(s))
}

func AtoInt32(s string) int32 {
	return int32(AToF64(s))
}

func AToF64Trunc(s string, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((AToF64(s)+0.5/n10)*n10) / n10
}

func AToF64(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return v
	}
	return 0
}

func AToU64(s string) uint64 {
	v, err := strconv.ParseUint(s, 10, 64)
	if err == nil {
		return v
	}
	return 0
}

func ToStr(s interface{}) string {
	return fmt.Sprintf("%v", s)
}

func StrToStrVec(str string, sep string) (vec []string) {
	for _, fs := range strings.Split(str, sep) {
		vec = append(vec, fs)
	}
	return
}

func StrToUintVec(str string, sep string) (vec []uint32) {
	for _, fs := range strings.Split(str, sep) {
		if f, err := strconv.Atoi(fs); err == nil {
			vec = append(vec, uint32(f))
		}
	}
	return
}
