package vm

import (
	"errors"
	"fmt"
	"os"

	"github.com/grubby/grubby/ast"
	"github.com/grubby/grubby/interpreter/vm/builtins"
	builtin_errors "github.com/grubby/grubby/interpreter/vm/builtins/errors"
	"github.com/grubby/grubby/parser"
)

type vm struct {
	ObjectSpace map[string]builtins.Value
}

type VM interface {
	Run(string) (builtins.Value, error)
	Get(string) (builtins.Value, error)
	MustGet(string) builtins.Value
}

func NewVM() VM {
	vm := &vm{
		ObjectSpace: make(map[string]builtins.Value),
	}
	objectClass := builtins.NewGlobalObjectClass()
	vm.ObjectSpace["Object"] = objectClass
	vm.ObjectSpace["Kernel"] = builtins.NewGlobalKernelClass()

	main := objectClass.(builtins.Class).New()
	main.AddMethod(builtins.NewMethod("to_s", func(args ...builtins.Value) (builtins.Value, error) {
		return builtins.NewString("main"), nil
	}))
	main.AddMethod(builtins.NewMethod("require", func(args ...builtins.Value) (builtins.Value, error) {
		fileName := args[0].(*builtins.StringValue)
		_, err := os.Open(fileName.String())
		if err != nil {
			err := fmt.Sprintf("LoadError: cannot load such file -- %s", fileName.String())
			return nil, builtin_errors.NewLoadError(err)
		}

		return nil, nil
	}))
	vm.ObjectSpace["main"] = main

	return vm
}

func (vm *vm) MustGet(key string) builtins.Value {
	return vm.ObjectSpace[key]
}

func (vm *vm) Get(key string) (builtins.Value, error) {
	val, ok := vm.ObjectSpace[key]
	if !ok {
		return nil, errors.New(fmt.Sprintf("NameError: undefined local variable or method '%s' for main:Object", key))
	}

	return val, nil
}

func (vm *vm) Run(input string) (builtins.Value, error) {
	lexer := parser.NewLexer(input)
	parser.RubyParse(lexer)

	main := vm.ObjectSpace["main"]
	return vm.executeWithContext(parser.Statements, main)
}

func (vm *vm) executeWithContext(statements []ast.Node, context builtins.Value) (builtins.Value, error) {
	var val builtins.Value
	var err error

	for _, statement := range statements {
		switch statement.(type) {
		case ast.FuncDecl:
			// FIXME: assumes for now this will only ever be at the top level
			// it seems like this should be replaced with context, but that's really context for calling methods, not
			// necessarily for defining new methods
			funcNode := statement.(ast.FuncDecl)
			method := builtins.NewMethod(funcNode.Name.Name, func(args ...builtins.Value) (builtins.Value, error) {
				return nil, nil
			})
			val = method
			vm.ObjectSpace["Kernel"].AddPrivateMethod(method)
		case ast.SimpleString:
			val = builtins.NewString(statement.(ast.SimpleString).Value)
		case ast.ConstantInt:
			val = builtins.NewInt(statement.(ast.ConstantInt).Value)
		case ast.CallExpression:
			var method builtins.Method
			callExpr := statement.(ast.CallExpression)

			// FIXME: this should not be a linear time lookup
			for _, m := range context.Methods() {
				if m.Name() == callExpr.Func.Name {
					method = m
					break
				}
			}

			if method == nil {
				err := builtin_errors.NewNameError(callExpr.Func.Name, context.String(), context.Class().String())
				return nil, err
			}

			args := []builtins.Value{}
			for _, astArgument := range callExpr.Args {
				arg, err := vm.executeWithContext([]ast.Node{astArgument}, context)
				if err != nil {
					return nil, err
				}

				args = append(args, arg)
			}

			val, err = method.Execute(args...)

		default:
			panic(fmt.Sprintf("handled unknown statement type: %T:\n\t\n => %#v\n", statement, statement))
		}
	}

	return val, err
}