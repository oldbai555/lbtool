package utils

import (
	"fmt"
	"testing"
)

func TestPluckStructField2IntList(t *testing.T) {
	type Str struct {
		Age int `json:"age"`
	}
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
	val := PluckStructField2IntList(&strList, "Age")
	fmt.Println(val)
}
