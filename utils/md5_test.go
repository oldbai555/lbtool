package utils

import (
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {
	md5 := Md5("hello world")
	fmt.Printf("%s", md5)
}
