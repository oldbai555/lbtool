package utils

import (
	"fmt"
	"testing"
)

func TestCamel2UnderScore(t *testing.T) {
	fmt.Println(Camel2UnderScore("TaaAAA"))
}
