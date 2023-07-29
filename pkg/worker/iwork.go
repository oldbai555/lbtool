package worker

// 在计算机编程术语里handle作为名词时是对可进行管理的资源对象的抽象，
// handle指向某个类别的资源对象，如文件句柄，进程ID都可以用handle来表达，在当动词讲时含义是处理和操作。
// 而handler表示的是过程（函数），理解为功能处理器的含义，如常用的回调函数可以用handler来表示。

type IWorker interface {
	Start(func(msg IMsg))
	Send(Type, interface{}) error
}

type IMsg interface {
	GetType() Type
	GetValue() interface{}
}
