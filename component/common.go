package component

import (
	"codnect.io/reflector"
	"fmt"
)

func canConvert(sourceType reflector.Type, targetType reflector.Type) bool {
	if sourceType.ReflectType() == targetType.ReflectType() {
		return true
	} else if reflector.IsInterface(targetType) && sourceType.CanConvert(targetType) {
		return true
	}

	return false
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
