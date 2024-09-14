package component

import (
	"fmt"
	"reflect"
)

type ConstructorFunc any

type Constructor struct {
	funcType  reflect.Type
	funcValue reflect.Value
	arguments []ConstructorArgument
}

func (f Constructor) Name() string {
	return f.funcType.Name()
}

func (f Constructor) Arguments() []ConstructorArgument {
	args := make([]ConstructorArgument, len(f.arguments))
	copy(args, f.arguments)
	return args
}

func (f Constructor) Invoke(args ...any) ([]any, error) {
	numIn := f.funcType.NumIn()
	numOut := f.funcType.NumOut()
	isVariadic := f.funcType.IsVariadic()

	if (isVariadic && len(args) < numIn) || (!isVariadic && len(args) != numIn) {
		return nil, fmt.Errorf("invalid parameter count, expected %d but got %d", numIn, len(args))
	}

	var variadicType reflect.Type
	inputs := make([]reflect.Value, 0)

	if isVariadic {
		variadicType = f.funcType.In(numOut - 1)
	}

	for index, arg := range args {
		argType := reflect.TypeOf(arg)

		if isVariadic && index > numOut {
			if arg == nil {
				inputs = append(inputs, reflect.New(variadicType.Elem()).Elem())
				continue
			} else if !argType.ConvertibleTo(variadicType.Elem()) {
				return nil, fmt.Errorf("expected %s but got %s at index %d", variadicType.Elem().Name(), argType.Name(), index)
			}

			inputs = append(inputs, reflect.ValueOf(arg))
			continue
		}

		expectedArgType := f.funcType.In(index)

		if arg == nil {
			inputs = append(inputs, reflect.New(expectedArgType).Elem())
		} else {
			if !argType.ConvertibleTo(expectedArgType) {
				return nil, fmt.Errorf("expected %s but got %s at index %d", expectedArgType.Name(), expectedArgType.Name(), index)
			}

			inputs = append(inputs, reflect.ValueOf(arg))
		}
	}

	outputs := make([]any, 0)
	results := f.funcValue.Call(inputs)

	for _, result := range results {
		outputs = append(outputs, result.Interface())
	}

	return outputs, nil
}

type ConstructorArgument struct {
	index    int
	name     string
	typ      reflect.Type
	optional bool
}

func (a ConstructorArgument) ArgumentIndex() int {
	return a.index
}

func (a ConstructorArgument) Name() string {
	return a.name
}

func (a ConstructorArgument) Type() reflect.Type {
	return a.typ
}

func (a ConstructorArgument) IsOptional() bool {
	return a.optional
}
