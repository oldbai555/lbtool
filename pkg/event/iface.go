package event

type Type uint32
type Fn func(IMsg) error
type Fns []Fn

type IEvent interface {
	Reg(t Type, fn Fn)
	Trigger(t Type, m IMsg)
}

type IMsg interface {
	GetValue() interface{}
}
