package filter

import (
	"reflect"
)

type Filter func(filters *Filters)

type Filters struct {
	Name string
	Type reflect.Type
}

func Of(filters ...Filter) *Filters {
	filterOpts := &Filters{}

	for _, filter := range filters {
		filter(filterOpts)
	}

	return filterOpts
}

func ByName(name string) Filter {
	return func(filters *Filters) {
		filters.Name = name
	}
}

func ByTypeOf[T any]() Filter {
	return func(filters *Filters) {
		typ := reflect.TypeFor[T]()
		filters.Type = typ
	}
}

func ByType(typ reflect.Type) Filter {
	return func(filters *Filters) {
		filters.Type = typ
	}
}
