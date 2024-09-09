package condition

import (
	"codnect.io/procyon-core/component"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

type AnyType struct {
}

func anyConstructorFunction() AnyType {
	return AnyType{}
}

func TestOnTypeCondition_MatchesConditionShouldReturnTrueIfAnyObjectWithTypeExists(t *testing.T) {
	onTypeCondition := OnType[AnyType]()
	container := component.NewObjectContainer()
	err := container.Singletons().Register("anyObject", AnyType{})
	assert.Nil(t, err)

	conditionContext := component.NewConditionContext(context.Background(), container)
	assert.True(t, onTypeCondition.MatchesCondition(conditionContext))
}

func TestOnTypeCondition_MatchesConditionShouldReturnTrueIfAnyDefinitionWithTypeExists(t *testing.T) {
	onTypeCondition := OnType[AnyType]()
	container := component.NewObjectContainer()

	definition, err := component.MakeDefinition(anyConstructorFunction)
	assert.Nil(t, err)
	err = container.Definitions().Register(definition)
	assert.Nil(t, err)

	conditionContext := component.NewConditionContext(context.Background(), container)
	assert.True(t, onTypeCondition.MatchesCondition(conditionContext))
}

func TestOnTypeCondition_MatchesConditionShouldReturnFalseIfAnyObjectWithTypeDoesNotExist(t *testing.T) {
	onTypeCondition := OnType[AnyType]()
	container := component.NewObjectContainer()

	conditionContext := component.NewConditionContext(context.Background(), container)
	assert.False(t, onTypeCondition.MatchesCondition(conditionContext))
}

func TestOnMissingTypeCondition_MatchesConditionShouldReturnFalseIfAnyObjectWithTypeExists(t *testing.T) {
	onMissingTypeCondition := OnMissingType[AnyType]()
	container := component.NewObjectContainer()
	err := container.Singletons().Register("anyObject", AnyType{})
	assert.Nil(t, err)

	conditionContext := component.NewConditionContext(context.Background(), container)
	assert.False(t, onMissingTypeCondition.MatchesCondition(conditionContext))
}

func TestOnMissingTypeCondition_MatchesConditionShouldReturnFalseIfAnyDefinitionWithTypeExists(t *testing.T) {
	onMissingTypeCondition := OnMissingType[AnyType]()
	container := component.NewObjectContainer()

	definition, err := component.MakeDefinition(anyConstructorFunction)
	assert.Nil(t, err)
	err = container.Definitions().Register(definition)
	assert.Nil(t, err)

	conditionContext := component.NewConditionContext(context.Background(), container)
	assert.False(t, onMissingTypeCondition.MatchesCondition(conditionContext))
}

func TestOnMissingTypeCondition_MatchesConditionShouldReturnTrueIfAnyObjectWithTypeDoesNotExist(t *testing.T) {
	onMissingTypeCondition := OnMissingObject("anyObjectName")
	container := component.NewObjectContainer()

	conditionContext := component.NewConditionContext(context.Background(), container)
	assert.True(t, onMissingTypeCondition.MatchesCondition(conditionContext))
}
