package builtins

type Array struct {
	valueStub
	members []Value
}

func NewArrayClass() Value {
	a := &Array{}
	a.initialize()
	return a
}

func (array *Array) String() string {
	return "Array"
}

func (array *Array) Append(v Value) {
	array.members = append(array.members, v)
}

func (array *Array) Members() []Value {
	return array.members
}