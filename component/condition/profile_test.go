package condition

import (
	"codnect.io/procyon-core/component"
	"codnect.io/procyon-core/runtime"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOnProfileCondition_MatchesConditionShouldReturnTrueIfProfileIsActivated(t *testing.T) {
	onProfileCondition := OnProfile("anyProfileName")
	container := component.NewContainer()

	environment := runtime.NewDefaultEnvironment()
	environment.SetActiveProfiles("anyProfileName")
	container.Singletons().Register("environment", environment)

	conditionContext := NewContext(context.Background(), container)
	assert.True(t, onProfileCondition.MatchesCondition(conditionContext))
}

func TestOnProfileCondition_MatchesConditionShouldReturnFalseIfEnvironmentObjectDoesNotExist(t *testing.T) {
	onProfileCondition := OnProfile("anyProfileName")
	container := component.NewContainer()

	conditionContext := NewContext(context.Background(), container)
	assert.False(t, onProfileCondition.MatchesCondition(conditionContext))
}

func TestOnProfileCondition_MatchesConditionShouldReturnFalseIfProfileIsNotActivated(t *testing.T) {
	onProfileCondition := OnProfile("anyProfileName")
	container := component.NewContainer()

	environment := runtime.NewDefaultEnvironment()
	environment.SetActiveProfiles("anotherProfileName")
	container.Singletons().Register("environment", environment)

	conditionContext := NewContext(context.Background(), container)
	assert.False(t, onProfileCondition.MatchesCondition(conditionContext))
}
