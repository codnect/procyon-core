package component

import (
	"codnect.io/reflector"
	"fmt"
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
		typ := reflector.TypeOf[T]()

		exists := false
		for _, arg := range definition.constructorArgs {
			if arg.Type().Compare(typ) {
				arg.name = name
				exists = true
			}
		}

		if !exists {
			return fmt.Errorf("could not find any input of type %s", typ.Name())
		}

		return nil
	}
}

func QualifierAt(index int, name string) Option {
	return func(definition *Definition) error {
		if index < 0 {
			panic(fmt.Sprintf("index should be greater than or equal to zero, but got index %d", index))
		}

		if len(definition.constructorArgs) <= index {
			return fmt.Errorf("could not find any input at index %d", index)
		}

		definition.constructorArgs[index].name = name
		return nil
	}
}

func Optional[T any]() Option {
	return func(def *Definition) error {
		typ := reflector.TypeOf[T]()

		exists := false
		for _, arg := range def.constructorArgs {
			if arg.Type().Compare(typ) {
				arg.optional = true
				exists = true
			}
		}

		if !exists {
			return fmt.Errorf("could not find any input of type %s", typ.Name())
		}

		return nil
	}
}

func OptionalAt(index int) Option {
	return func(definition *Definition) error {
		if index < 0 {
			panic(fmt.Sprintf("index should be greater than or equal to zero, but got index %d", index))
		}

		if len(definition.constructorArgs) <= index {
			return fmt.Errorf("could not find any input at index %d", index)
		}

		definition.constructorArgs[index].optional = true
		return nil
	}
}
