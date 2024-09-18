package condition

import (
	"codnect.io/procyon-core/component"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOnObjectCondition_MatchesConditionShouldReturnTrueIfAnyObjectWithNameExists(t *testing.T) {
	onObjectCondition := OnObject("anyObjectName")
	container := component.NewContainer()
	err := container.Singletons().Register("anyObjectName", "anyObject")
	assert.Nil(t, err)

	conditionContext := NewContext(context.Background(), container)
	assert.True(t, onObjectCondition.MatchesCondition(conditionContext))
}

func TestOnObjectCondition_MatchesConditionShouldReturnTrueIfAnyDefinitionWithNameExists(t *testing.T) {
	onObjectCondition := OnObject("anyObjectName")
	container := component.NewContainer()

	definition, err := component.MakeDefinition(anyConstructorFunction, component.Named("anyObjectName"))
	assert.Nil(t, err)
	err = container.Definitions().Register(definition)
	assert.Nil(t, err)

	conditionContext := NewContext(context.Background(), container)
	assert.True(t, onObjectCondition.MatchesCondition(conditionContext))
}

func TestOnObjectCondition_MatchesConditionShouldReturnFalseIfAnyObjectWithNameDoesNotExist(t *testing.T) {
	onObjectCondition := OnObject("anyObjectName")
	container := component.NewContainer()

	conditionContext := NewContext(context.Background(), container)
	assert.False(t, onObjectCondition.MatchesCondition(conditionContext))
}

func TestOnMissingObjectCondition_MatchesConditionShouldReturnFalseIfAnyObjectWithNameExists(t *testing.T) {
	onMissingObjectCondition := OnMissingObject("anyObjectName")
	container := component.NewContainer()
	err := container.Singletons().Register("anyObjectName", "anyObject")
	assert.Nil(t, err)

	conditionContext := NewContext(context.Background(), container)
	assert.False(t, onMissingObjectCondition.MatchesCondition(conditionContext))
}

func TestOnMissingObjectCondition_MatchesConditionShouldReturnFalseIfAnyDefinitionWithNameExists(t *testing.T) {
	onMissingObjectCondition := OnMissingObject("anyObjectName")
	container := component.NewContainer()

	definition, err := container.MakeDefinition(anyConstructorFunction, component.Named("anyObjectName"))
	assert.Nil(t, err)
	err = container.Definitions().Register(definition)
	assert.Nil(t, err)

	conditionContext := NewContext(context.Background(), container)
	assert.False(t, onMissingObjectCondition.MatchesCondition(conditionContext))
}

func TestOnMissingObjectCondition_MatchesConditionShouldReturnTrueIfAnyObjectWithNameDoesNotExist(t *testing.T) {
	onMissingObjectCondition := OnMissingObject("anyObjectName")
	container := component.NewContainer()

	conditionContext := NewContext(context.Background(), container)
	assert.True(t, onMissingObjectCondition.MatchesCondition(conditionContext))
}
