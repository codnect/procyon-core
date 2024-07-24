package filter

import "codnect.io/reflector"

type Filter func(filters *Filters)

type Filters struct {
	Name      string
	Type      reflector.Type
	Arguments []any
}

func ByName(name string) Filter {
	return func(filters *Filters) {
		filters.Name = name
	}
}

func ByTypeOf[T any]() Filter {
	return func(filters *Filters) {
		typ := reflector.TypeOf[T]()
		filters.Type = typ
	}
}

func ByType(typ reflector.Type) Filter {
	return func(filters *Filters) {
		filters.Type = typ
	}
}

func ByArguments(args ...any) Filter {
	return func(filters *Filters) {
		if len(args) != 0 {
			filters.Arguments = append(filters.Arguments, args...)
		}
	}
}
