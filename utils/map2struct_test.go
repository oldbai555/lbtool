package utils

import (
	"log"
	"testing"
)

type TestStructure struct {
	Field1 string
	Field2 string
}

func TestMap2Struct(t *testing.T) {
	var valMap = make(map[string]interface{})
	valMap["Field1"] = "value1"
	valMap["Field2"] = "value2"
	var testStructure = TestStructure{}
	err := Map2Struct(valMap, &testStructure)
	if err != nil {
		log.Println(err)
	}
	log.Println(testStructure)
}
