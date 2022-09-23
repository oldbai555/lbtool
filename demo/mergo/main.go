package main

import (
	"fmt"
	"github.com/imdario/mergo"
	"log"
)

type Student struct {
	MyName string
	Num    int
	Age    int
}

func main() {
	var defaultStudent = Student{
		MyName: "zhang—san",
		Num:    1,
		Age:    18,
	}

	var m = make(map[string]interface{})
	if err := mergo.Map(&m, defaultStudent); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("map m = %+v", m)

	var defaultStudent1 Student
	err := mergo.Map(&defaultStudent1, m)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("struct  s = %+v", defaultStudent1)
}
