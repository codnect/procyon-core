package component

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestObjectNotFoundError_ErrorShouldReturnErrorMessageIfNameIsProvided(t *testing.T) {
	err := ObjectNotFoundError{
		name: "anyObjectName",
	}

	assert.Equal(t, "not found object with name 'anyObjectName'", err.Error())
}

func TestObjectNotFoundError_ErrorShouldReturnErrorMessageIfTypeIsProvided(t *testing.T) {
	err := ObjectNotFoundError{
		typ: reflect.TypeFor[AnyType](),
	}

	assert.Equal(t, "not found object with type 'component.AnyType'", err.Error())
}

func TestObjectNotFoundError_ErrorShouldReturnErrorMessageIfTypeAndNameAreProvided(t *testing.T) {
	err := ObjectNotFoundError{
		name: "anyObjectName",
		typ:  reflect.TypeFor[AnyType](),
	}

	assert.Equal(t, "not found object with name 'anyObjectName' and type 'component.AnyType'", err.Error())
}

func TestDefinitionNotFoundError_ErrorShouldReturnErrorMessageIfNameIsProvided(t *testing.T) {
	err := DefinitionNotFoundError{
		name: "anyObjectName",
	}

	assert.Equal(t, "not found definition with name 'anyObjectName'", err.Error())
}

func TestDefinitionNotFoundError_ErrorShouldReturnErrorMessageIfTypeIsProvided(t *testing.T) {
	err := DefinitionNotFoundError{
		typ: reflect.TypeFor[AnyType](),
	}

	assert.Equal(t, "not found definition with type 'component.AnyType'", err.Error())
}

func TestDefinitionNotFoundError_ErrorShouldReturnErrorMessageIfTypeAndNameAreProvided(t *testing.T) {
	err := DefinitionNotFoundError{
		name: "anyObjectName",
		typ:  reflect.TypeFor[AnyType](),
	}

	assert.Equal(t, "not found definition with name 'anyObjectName' and type 'component.AnyType'", err.Error())
}
