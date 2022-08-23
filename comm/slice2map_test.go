package comm

import (
	"fmt"
	"testing"
)

type Str struct {
	Age int `json:"age"`
}

func TestSliceToMapKeyByStructField(t *testing.T) {
	var strList = []*Str{
		{
			Age: 1,
		},
		{
			Age: 2,
		},
		{
			Age: 3,
		},
	}
	val := Slice2MapKeyByStructField(&strList, "Age").(map[int]*Str)
	fmt.Println(val)
}

func TestSliceStruct2MapValueByBool(t *testing.T) {
	var strList = []*Str{
		{
			Age: 1,
		},
		{
			Age: 2,
		},
		{
			Age: 3,
		},
	}
	val := SliceStruct2MapKeyFieldValueByBool(&strList, "Age").(map[int]bool)
	fmt.Println(val[4])
	fmt.Println(val[3])
}

func TestSliceBasis2MapValueByBool(t *testing.T) {
	var val = []string{"1", "2", "3", "4", "5", "6", "7"}
	res := SliceBasis2MapValueByBool(val).(map[string]bool)
	fmt.Println(res)
}
