package component

import (
	"codnect.io/reflector"
	"github.com/stretchr/testify/assert"
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
		typ: reflector.TypeOf[AnyType](),
	}

	assert.Equal(t, "not found object with type 'AnyType'", err.Error())
}

func TestObjectNotFoundError_ErrorShouldReturnErrorMessageIfTypeAndNameAreProvided(t *testing.T) {
	err := ObjectNotFoundError{
		name: "anyObjectName",
		typ:  reflector.TypeOf[AnyType](),
	}

	assert.Equal(t, "not found object with name 'anyObjectName' and type 'AnyType'", err.Error())
}

func TestDefinitionNotFoundError_ErrorShouldReturnErrorMessageIfNameIsProvided(t *testing.T) {
	err := DefinitionNotFoundError{
		name: "anyObjectName",
	}

	assert.Equal(t, "not found definition with name 'anyObjectName'", err.Error())
}

func TestDefinitionNotFoundError_ErrorShouldReturnErrorMessageIfTypeIsProvided(t *testing.T) {
	err := DefinitionNotFoundError{
		typ: reflector.TypeOf[AnyType](),
	}

	assert.Equal(t, "not found definition with type 'AnyType'", err.Error())
}

func TestDefinitionNotFoundError_ErrorShouldReturnErrorMessageIfTypeAndNameAreProvided(t *testing.T) {
	err := DefinitionNotFoundError{
		name: "anyObjectName",
		typ:  reflector.TypeOf[AnyType](),
	}

	assert.Equal(t, "not found definition with name 'anyObjectName' and type 'AnyType'", err.Error())
}
