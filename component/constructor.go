package component

import "codnect.io/reflector"

type Constructor any

type ConstructorArgument struct {
	index    int
	name     string
	typ      reflector.Type
	optional bool
}

func (a ConstructorArgument) ArgumentIndex() int {
	return a.index
}

func (a ConstructorArgument) Name() string {
	return a.name
}

func (a ConstructorArgument) Type() reflector.Type {
	return a.typ
}

func (a ConstructorArgument) IsOptional() bool {
	return a.optional
}
