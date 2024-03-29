package container

import (
	"codnect.io/reflector"
	"fmt"
)

type Definition struct {
	name            string
	typ             reflector.Type
	constructorFunc reflector.Function
	constructor     Constructor
	scope           string
	inputs          []*Input
}

func (d *Definition) Name() string {
	return d.name
}

func (d *Definition) Type() reflector.Type {
	return d.typ
}

func (d *Definition) Constructor() reflector.Function {
	return d.constructorFunc
}

func (d *Definition) Scope() string {
	return d.scope
}

func (d *Definition) Inputs() []*Input {
	inputs := make([]*Input, 0)

	for _, input := range d.inputs {
		inputs = append(inputs, input)
	}

	return inputs
}

func (d *Definition) IsShared() bool {
	return d.scope == SharedScope
}

func (d *Definition) IsPrototype() bool {
	return d.scope == PrototypeScope
}

func MakeDefinition(constructor Constructor, options ...Option) (*Definition, error) {
	if constructor == nil {
		return nil, fmt.Errorf("container: constructor should not be nil")
	}

	typ := reflector.TypeOfAny(constructor)
	if !reflector.IsFunction(typ) {
		return nil, fmt.Errorf("container: constructor should be a function")
	}

	functionType := reflector.ToFunction(typ)

	if functionType.NumResult() != 1 {
		return nil, fmt.Errorf("container: constructor can only be a function returning one result")
	}

	resultType := functionType.Results()[0]

	name := ""

	if reflector.IsPointer(resultType) {
		pointerType := reflector.ToPointer(resultType)
		name = pointerType.Elem().Name()
	} else {
		name = resultType.Name()
	}

	def := &Definition{
		name:            lowerCamelCase(name),
		constructor:     constructor,
		constructorFunc: functionType,
		typ:             resultType,
		scope:           SharedScope,
	}

	for index, parameterType := range functionType.Parameters() {
		input := &Input{
			index: index,
			typ:   parameterType,
		}

		def.inputs = append(def.inputs, input)
	}

	for _, option := range options {
		err := option(def)
		if err != nil {
			return nil, err
		}
	}

	return def, nil
}
