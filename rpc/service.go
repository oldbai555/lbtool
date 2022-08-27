package rpc

import (
	"fmt"
	"github.com/oldbai555/lb/log"
	"go/ast"
	"reflect"
	"sync/atomic"
)

type service struct {
	// name 即映射的结构体的名称;
	name string
	// typ 是结构体的类型;
	typ reflect.Type
	// rcvr 即结构体的实例本身，保留 rcvr 是因为在调用时需要 rcvr 作为第 0 个参数;
	rcvr reflect.Value
	// method 是 map 类型，存储映射的结构体的所有符合条件的方法;
	method map[string]*methodType
}

// newService 构造成指针结构体
func newService(rcvr interface{}) *service {
	s := new(service)
	s.rcvr = reflect.ValueOf(rcvr)
	s.name = reflect.Indirect(s.rcvr).Type().Name()
	s.typ = reflect.TypeOf(rcvr)
	if !ast.IsExported(s.name) {
		panic(any(fmt.Sprintf("rpc server: %s is not a valid service name", s.name)))
	}
	s.registerMethods()
	return s
}

// registerMethods 过滤出了符合条件的方法
func (s *service) registerMethods() {
	s.method = make(map[string]*methodType)

	for i := 0; i < s.typ.NumMethod(); i++ {
		method := s.typ.Method(i)
		mType := method.Type

		// 两个导出或内置类型的入参（反射时为 3 个，第 0 个是自身，类似于 python 的 self，java 中的 this）
		// 返回值有且只有 1 个，类型为 error
		if mType.NumIn() != 3 || mType.NumOut() != 1 {
			continue
		}

		if mType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}

		argType, replyType := mType.In(1), mType.In(2)
		if !isExportedOrBuiltinType(argType) || !isExportedOrBuiltinType(replyType) {
			continue
		}

		s.method[method.Name] = &methodType{
			method:    method,
			ArgType:   argType,
			ReplyType: replyType,
		}
		log.Infof("rpc server: register %s.%s", s.name, method.Name)
	}

}

// call 通过反射值调用方法
func (s *service) call(m *methodType, argv, replyv reflect.Value) error {
	atomic.AddUint64(&m.numCalls, 1)
	f := m.method.Func
	returnValues := f.Call([]reflect.Value{s.rcvr, argv, replyv})
	if errInter := returnValues[0].Interface(); errInter != nil {
		return errInter.(error)
	}
	return nil
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}
