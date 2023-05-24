package all_sort

import (
	"fmt"
	"strings"
	"testing"
)

func Test7Base(t *testing.T) {
	in := "1 2 3 4"
	result := outOrder(strings.Fields(in))
	// dictSort(result)
	// s := format(result)
	// fmt.Println(s)
	fmt.Println("len is ", len(result))
	for _, strArr := range result {
		fmt.Println(strings.Join(strArr, ","))
	}
}
