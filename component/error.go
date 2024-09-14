package component

import (
	"fmt"
	"reflect"
	"strings"
)

type ObjectNotFoundError struct {
	name string
	typ  reflect.Type
}

func (e ObjectNotFoundError) Error() string {
	var builder strings.Builder
	builder.WriteString("not found object with")

	if e.name != "" {
		builder.WriteString(fmt.Sprintf(" name '%s'", e.name))
		if e.typ != nil {
			builder.WriteString(" and")
		}
	}

	if e.typ != nil {
		builder.WriteString(fmt.Sprintf(" type '%s'", e.typ.String()))
	}

	return builder.String()
}

type DefinitionNotFoundError struct {
	name string
	typ  reflect.Type
}

func (e DefinitionNotFoundError) Error() string {
	var builder strings.Builder
	builder.WriteString("not found definition with")

	if e.name != "" {
		builder.WriteString(fmt.Sprintf(" name '%s'", e.name))
		if e.typ != nil {
			builder.WriteString(" and")
		}
	}

	if e.typ != nil {
		builder.WriteString(fmt.Sprintf(" type '%s'", e.typ.String()))
	}

	return builder.String()
}
