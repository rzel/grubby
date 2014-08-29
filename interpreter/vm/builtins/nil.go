package builtins

type NilClass struct {
	valueStub
}

func NewNilClass() Value {
	n := &NilClass{}
	n.initialize()
	n.class = NewClassValue().(Class)
	return n
}

func (n *NilClass) String() string {
	return "NilClass"
}

type nilInstance struct {
	valueStub
}

func Nil() Value {
	return NewNilClass().(Class).New()
}

func (class *NilClass) New() Value {
	n := &nilInstance{}
	n.initialize()
	n.class = class

	return n
}

func (n *nilInstance) String() string {
	return ""
}