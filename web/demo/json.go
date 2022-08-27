package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// golang反射创建对象，解析JSON
func main() {
	fv := reflect.ValueOf(Hello)
	ft := fv.Type()
	if fv.Kind() == reflect.Func {
		fmt.Println(fv.String())
		fmt.Println(ft.String())

		fmt.Println(ft.In(0).String())
		fmt.Println(ft.In(1).String())

		reqNewV := reflect.New(ft.In(1).Elem())
		fmt.Println(reqNewV.Type().String())
		msg := reqNewV.Interface().(TestInterface)

		var json = "{\n    \"name\": \"lili\"\n}"
		err := UnmarshalTestHello([]byte(json), msg)
		if err != nil {
			fmt.Printf("err:%v\n", err)
			return
		}

		var params []reflect.Value
		params = append(params, reflect.ValueOf("aaa"))
		params = append(params, reflect.ValueOf(msg))
		rs := fv.Call(params)

		fmt.Println("result:", rs[0].Interface().(string))
		fmt.Println("err:", rs[1].Interface())
	}
}

type TestInterface interface {
	Hello() string
}

func UnmarshalTestHello(data []byte, r TestInterface) error {
	err := json.Unmarshal(data, r)
	return err
}

type TestHello struct {
	Name string `json:"name"`
}

func (r *TestHello) Hello() string {
	panic(any("implement me"))
}

func Hello(val string, hello *TestHello) (string, error) {
	fmt.Println(hello.Name)
	return val, nil
}
