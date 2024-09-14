package component

import (
	"fmt"
	"reflect"
	"strings"
)

func Named(name string) Option {
	return func(definition *Definition) error {
		if strings.TrimSpace(name) != "" {
			definition.name = name
		}

		return nil
	}
}

func Scoped(scope string) Option {
	return func(definition *Definition) error {
		if strings.TrimSpace(scope) == "" {
			definition.scope = SingletonScope
		} else {
			definition.scope = scope
		}

		return nil
	}
}

func Qualifier[T any](name string) Option {
	return func(definition *Definition) error {
		typ := reflect.TypeFor[T]()
		constructor := definition.Constructor()

		exists := false
		for index, arg := range constructor.Arguments() {
			if arg.Type() == typ {
				constructor.arguments[index].name = name
				exists = true
			}
		}

		if !exists {
			return fmt.Errorf("cannot find any input of type %s", typ.Name())
		}

		return nil
	}
}

func QualifierAt(index int, name string) Option {
	return func(definition *Definition) error {
		if index < 0 {
			panic(fmt.Sprintf("index should be greater than or equal to zero, but got index %d", index))
		}

		constructor := definition.Constructor()
		if len(constructor.Arguments()) <= index {
			return fmt.Errorf("cannot find any input at index %d", index)
		}

		constructor.arguments[index].name = name
		return nil
	}
}

func Optional[T any]() Option {
	return func(definition *Definition) error {
		typ := reflect.TypeFor[T]()
		constructor := definition.Constructor()

		exists := false
		for index, arg := range constructor.Arguments() {
			if arg.Type() == typ {
				constructor.arguments[index].optional = true
				exists = true
			}
		}

		if !exists {
			return fmt.Errorf("cannot find any input of type %s", typ.Name())
		}

		return nil
	}
}

func OptionalAt(index int) Option {
	return func(definition *Definition) error {
		if index < 0 {
			panic(fmt.Sprintf("index should be greater than or equal to zero, but got index %d", index))
		}

		constructor := definition.Constructor()
		if len(constructor.Arguments()) <= index {
			return fmt.Errorf("cannot find any input at index %d", index)
		}

		constructor.arguments[index].optional = true
		return nil
	}
}
