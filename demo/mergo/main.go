package main

import (
	"fmt"
	"github.com/imdario/mergo"
	"log"
)

type Student struct {
	MyName string `json:"my_name,omitempty"`
	Num    int    `json:"num,omitempty"`
	Age    int    `json:"age,omitempty"`
	Cate   *Cate
}

type Cate struct {
	Id int64 `json:"id,omitempty"`
}

func main() {
	var defaultStudent = Student{
		MyName: "zhang—san",
		Age:    18,
		Cate: &Cate{
			Id: 1,
		},
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
