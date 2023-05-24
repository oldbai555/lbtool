package utils

import (
	"fmt"
	"sort"
	"testing"
)

type AStr struct {
	Age int `json:"age"`
}

type AStrList []*AStr

// Len 重写 Len() 方法
func (a AStrList) Len() int {
	return len(a)
}

// Swap 重写 Swap() 方法
func (a AStrList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less 重写 Less() 方法， 从大到小排序
func (a AStrList) Less(i, j int) bool {
	//return a[j].StartAt > a[i].StartAt // 从小到大
	return ReflectCompareFieldAsc(a[j], a[i], "Age") // 从大到小
}

func TestKindOfData(t *testing.T) {
	var aList AStrList
	aList = append(aList, &AStr{Age: 1})
	aList = append(aList, &AStr{Age: 3})
	aList = append(aList, &AStr{Age: 2})
	aList = append(aList, &AStr{Age: 4})
	aList = append(aList, &AStr{Age: 5})

	sort.Sort(aList)
	for _, str := range aList {
		fmt.Println(str)
	}
}
