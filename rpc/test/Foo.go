package test

import "fmt"

type Foo int

type Args struct{ Num1, Num2 int }

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

// it's not a exported Method
func (f Foo) sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func Assert(condition bool, msg string, v ...interface{}) {
	if !condition {
		val := fmt.Sprintf("assertion failed: %s,%v", msg, v)
		fmt.Println(val)
		panic(any(val))
	}
}
