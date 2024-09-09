package condition

import (
	"codnect.io/procyon-core/component"
	"codnect.io/procyon-core/runtime"
	"codnect.io/procyon-core/runtime/property"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOnPropertyCondition_MatchesConditionShouldReturnTrueIfPropertyExists(t *testing.T) {
	onPropertyCondition := OnProperty("anyPropertyName")
	container := component.NewObjectContainer()

	anyPropertySource := property.NewMapSource("anyPropertySource", map[string]interface{}{
		"anyPropertyName": true,
	})

	environment := runtime.NewDefaultEnvironment()
	environment.PropertySources().AddLast(anyPropertySource)
	container.Singletons().Register("environment", environment)

	conditionContext := component.NewConditionContext(context.Background(), container)
	assert.True(t, onPropertyCondition.MatchesCondition(conditionContext))
}

func TestOnPropertyCondition_MatchesConditionShouldReturnTrueEvenIfPropertyDoesNotExistAndMatchIfMissingIsCalled(t *testing.T) {
	onPropertyCondition := OnProperty("anyPropertyName").MatchIfMissing(true)
	container := component.NewObjectContainer()

	environment := runtime.NewDefaultEnvironment()
	container.Singletons().Register("environment", environment)

	conditionContext := component.NewConditionContext(context.Background(), container)
	assert.True(t, onPropertyCondition.MatchesCondition(conditionContext))
}

func TestOnPropertyCondition_MatchesConditionShouldReturnTrueIfPropertyValueEqualsToGivenValue(t *testing.T) {
	onPropertyCondition := OnProperty("anyPropertyName").HavingValue("anyPropertyValue")
	container := component.NewObjectContainer()

	anyPropertySource := property.NewMapSource("anyPropertySource", map[string]interface{}{
		"anyPropertyName": "anyPropertyValue",
	})

	environment := runtime.NewDefaultEnvironment()
	environment.PropertySources().AddLast(anyPropertySource)
	container.Singletons().Register("environment", environment)

	conditionContext := component.NewConditionContext(context.Background(), container)
	assert.True(t, onPropertyCondition.MatchesCondition(conditionContext))
}

func TestOnPropertyCondition_MatchesConditionShouldReturnTrueIfPropertyValueDoesNotEqualToGivenValue(t *testing.T) {
	onPropertyCondition := OnProperty("anyPropertyName").HavingValue("anotherPropertyValue")
	container := component.NewObjectContainer()

	anyPropertySource := property.NewMapSource("anyPropertySource", map[string]interface{}{
		"anyPropertyName": "anyPropertyValue",
	})

	environment := runtime.NewDefaultEnvironment()
	environment.PropertySources().AddLast(anyPropertySource)
	container.Singletons().Register("environment", environment)

	conditionContext := component.NewConditionContext(context.Background(), container)
	assert.False(t, onPropertyCondition.MatchesCondition(conditionContext))
}

func TestOnPropertyCondition_MatchesConditionShouldReturnFalseIfEnvironmentObjectDoesNotExist(t *testing.T) {
	onPropertyCondition := OnProperty("anyPropertyName")
	container := component.NewObjectContainer()

	conditionContext := component.NewConditionContext(context.Background(), container)
	assert.False(t, onPropertyCondition.MatchesCondition(conditionContext))
}
