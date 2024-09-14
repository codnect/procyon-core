package component

import (
	"reflect"
)

func matchTypes(sourceType reflect.Type, targetType reflect.Type) bool {
	if sourceType == targetType || (targetType.Kind() == reflect.Interface && sourceType.ConvertibleTo(targetType)) {
		return true
	} else if sourceType.Kind() == reflect.Pointer {
		return matchTypes(sourceType.Elem(), targetType)
	}

	return false
}
