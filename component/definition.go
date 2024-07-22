package component

import (
	"codnect.io/reflector"
	"fmt"
	"strings"
	"unicode"
)

type DefinitionOption func(definition *Definition)

type Definition struct {
	name            string
	typ             reflector.Type
	scope           string
	priority        int
	primary         bool
	constructor     reflector.Function
	constructorArgs []*ConstructorArgument
}

func MakeDefinition(constructor Constructor, options ...DefinitionOption) (*Definition, error) {
	if constructor == nil {
		return nil, fmt.Errorf("constructor should not be nil")
	}

	constructorType := reflector.TypeOfAny(constructor)
	if !reflector.IsFunction(constructorType) {
		return nil, fmt.Errorf("constructor should be a function")
	}

	constructorFunc := reflector.ToFunction(constructorType)
	if constructorFunc.NumResult() != 1 {
		return nil, fmt.Errorf("constructor can only be a function returning one result")
	}

	returnType := constructorFunc.Results()[0]
	definitionName := ""

	if reflector.IsPointer(returnType) {
		pointerType := reflector.ToPointer(returnType)
		definitionName = pointerType.Elem().Name()
	} else {
		definitionName = returnType.Name()
	}

	definition := &Definition{
		name:            lowerCamelCase(definitionName),
		typ:             returnType,
		scope:           SingletonScope,
		constructor:     constructorFunc,
		constructorArgs: make([]*ConstructorArgument, 0),
	}

	for index, parameterType := range constructorFunc.Parameters() {
		arg := &ConstructorArgument{
			index:        index,
			typ:          parameterType,
			requiredType: parameterType,
			optional:     false,
		}

		definition.constructorArgs = append(definition.constructorArgs, arg)
	}

	for _, option := range options {
		option(definition)
	}

	return definition, nil
}

func (d *Definition) Name() string {
	return d.name
}

func (d *Definition) Type() reflector.Type {
	return d.typ
}

func (d *Definition) Scope() string {
	return d.scope
}

func (d *Definition) Priority() int {
	return d.priority
}

func (d *Definition) IsPrimary() bool {
	return d.primary
}

func (d *Definition) IsSingleton() bool {
	return d.scope == SingletonScope
}

func (d *Definition) IsPrototype() bool {
	return d.scope == PrototypeScope
}

func (d *Definition) ConstructorArguments() []*ConstructorArgument {
	copyOfArgs := make([]*ConstructorArgument, 0)

	for _, arg := range d.constructorArgs {
		copyOfArgs = append(copyOfArgs, arg)
	}

	return copyOfArgs
}

func WithName(name string) DefinitionOption {
	return func(definition *Definition) {
		if strings.TrimSpace(name) != "" {
			definition.name = name
		}
	}
}

func WithPriority(priority int) DefinitionOption {
	return func(definition *Definition) {
		definition.priority = priority
	}
}

func WithPrimary() DefinitionOption {
	return func(definition *Definition) {
		definition.primary = true
	}
}

func WithScope(scope string) DefinitionOption {
	return func(definition *Definition) {
		if strings.TrimSpace(scope) == "" {
			definition.scope = SingletonScope
		} else {
			definition.scope = scope
		}
	}
}

func WithNamedArgument(index int, name string) DefinitionOption {
	return func(definition *Definition) {
		definition.constructorArgs[index].name = name
	}
}

func WithTypedArgument[T any](index int) DefinitionOption {
	return func(definition *Definition) {
		definition.constructorArgs[index].requiredType = reflector.TypeOf[T]()
	}
}

func WithOptionalArgument(index int) DefinitionOption {
	return func(definition *Definition) {
		definition.constructorArgs[index].optional = true
	}
}

func lowerCamelCase(str string) string {
	isFirst := true

	return strings.Map(func(r rune) rune {
		if isFirst {
			isFirst = false
			return unicode.ToLower(r)
		}

		return r
	}, str)

}
