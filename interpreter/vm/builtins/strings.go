package builtins

import (
	"fmt"
	"strconv"
)

type StringClass struct {
	valueStub
	classStub

	provider ClassProvider
}

func NewStringClass(classProvider ClassProvider, singletonProvider SingletonProvider) Class {
	s := &StringClass{}
	s.initialize()
	s.setStringer(s.String)

	s.provider = classProvider
	s.class = classProvider.ClassWithName("Class")
	s.superClass = classProvider.ClassWithName("Object")

	s.AddMethod(NewNativeMethod("+", classProvider, singletonProvider, func(self Value, block Block, args ...Value) (Value, error) {
		arg := args[0].(*StringValue)
		selfAsStr := self.(*StringValue)
		return NewString(selfAsStr.value+arg.value, classProvider, singletonProvider), nil
	}))
	s.AddMethod(NewNativeMethod("==", classProvider, singletonProvider, func(self Value, block Block, args ...Value) (Value, error) {
		asStr, ok := args[0].(*StringValue)
		if !ok {
			return singletonProvider.SingletonWithName("false"), nil
		}

		selfAsStr := self.(*StringValue)
		if selfAsStr.value == asStr.value {
			return singletonProvider.SingletonWithName("true"), nil
		} else {
			return singletonProvider.SingletonWithName("false"), nil
		}
	}))
	s.AddMethod(NewNativeMethod("to_i", classProvider, singletonProvider, func(self Value, block Block, args ...Value) (Value, error) {
		selfAsStr := self.(*StringValue)

		intValue, _ := strconv.Atoi(selfAsStr.value)
		return NewFixnum(intValue, classProvider, singletonProvider), nil
	}))

	return s
}

func (c *StringClass) String() string {
	return "String"
}

func (c *StringClass) Name() string {
	return "String"
}

func (class *StringClass) New(classProvider ClassProvider, singletonProvider SingletonProvider, args ...Value) (Value, error) {
	str := &StringValue{}
	str.initialize()
	str.setStringer(str.String)
	str.setStringer(str.String)
	str.class = class

	return str, nil
}

type StringValue struct {
	value string
	valueStub
}

func (s *StringValue) String() string {
	return fmt.Sprintf(`"%s"`, s.value)
}

func (s *StringValue) RawString() string {
	return s.value
}

func NewString(str string, classProvider ClassProvider, singletonProvider SingletonProvider) Value {
	s, _ := classProvider.ClassWithName("String").New(classProvider, singletonProvider)
	s.(*StringValue).value = str
	return s
}
