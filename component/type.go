package component

import (
	"codnect.io/reflector"
	"fmt"
)

func matchTypeName(sourceType reflector.Type, targetType reflector.Type) bool {
	return fullTypeName(sourceType) == fullTypeName(targetType)
}

func rawName(typ reflector.Type) string {
	if reflector.IsPointer(typ) {
		pointerTyp := reflector.ToPointer(typ)
		typ = pointerTyp.Elem()
	}

	return typ.Name()
}

func fullTypeName(typ reflector.Type) string {
	if reflector.IsPointer(typ) {
		pointerTyp := reflector.ToPointer(typ)
		typ = pointerTyp.Elem()
	}

	return fmt.Sprintf("%s.%s", typ.PackagePath(), typ.Name())
}
