package component

import (
	"codnect.io/procyon-core/component/filter"
	"context"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestDefinitionLoader_LoadDefinitionsShouldLoadDefinitionsIfComponentConditionsAreMet(t *testing.T) {
	container := NewObjectContainer()
	loader := NewDefinitionLoader(container)

	ctx := context.Background()
	conditionContext := NewConditionContext(ctx, container)

	mockCondition := MockCondition{}
	mockCondition.On("MatchesCondition", conditionContext).Return(true)

	component := createComponent(anyConstructorFunction, Named("anyObjectName"))
	component.conditions = append(component.conditions, mockCondition)

	err := loader.LoadDefinitions(ctx, []*Component{component})
	assert.Nil(t, err)

	var definition *Definition
	definition, err = container.Definitions().Find(filter.ByName("anyObjectName"))
	assert.Nil(t, err)
	assert.NotNil(t, definition)

	assert.Equal(t, "anyObjectName", definition.Name())
	assert.Equal(t, reflect.TypeFor[*AnyType](), definition.Type())
	assert.Equal(t, SingletonScope, definition.Scope())
	assert.True(t, definition.IsSingleton())
	assert.False(t, definition.IsPrototype())
	assert.NotNil(t, definition.Constructor())
	assert.Len(t, definition.Constructor().Arguments(), 0)
}

func TestDefinitionLoader_LoadDefinitionsShouldNotLoadDefinitionsIfComponentConditionsAreNotMet(t *testing.T) {
	container := NewObjectContainer()
	loader := NewDefinitionLoader(container)

	ctx := context.Background()
	conditionContext := NewConditionContext(ctx, container)

	mockCondition := MockCondition{}
	mockCondition.On("MatchesCondition", conditionContext).Return(false)

	component := createComponent(anyConstructorFunction, Named("anyObjectName"))
	component.conditions = append(component.conditions, mockCondition)

	err := loader.LoadDefinitions(ctx, []*Component{component})
	assert.Nil(t, err)

	var definition *Definition
	definition, err = container.Definitions().Find(filter.ByName("anyObjectName"))
	assert.Equal(t, "not found definition with name 'anyObjectName'", err.Error())
	assert.Nil(t, definition)
}

func TestDefinitionLoader_LoadDefinitionsShouldReturnErrorInCaseOfDuplicatedComponents(t *testing.T) {
	container := NewObjectContainer()
	loader := NewDefinitionLoader(container)

	ctx := context.Background()
	conditionContext := NewConditionContext(ctx, container)

	mockCondition := MockCondition{}
	mockCondition.On("MatchesCondition", conditionContext).Return(true)

	component := createComponent(anyConstructorFunction, Named("anyObjectName"))
	component.conditions = append(component.conditions, mockCondition)

	err := loader.LoadDefinitions(ctx, []*Component{component, component})
	assert.Equal(t, "definition with name 'anyObjectName' already exists", err.Error())

}
