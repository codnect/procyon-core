package condition

import (
	"codnect.io/procyon-core/component"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockCondition struct {
	mock.Mock
}

func (t MockCondition) MatchesCondition(ctx Context) bool {
	args := t.Called(ctx)
	return args.Bool(0)
}

func TestNewConditionContextShouldCreateContextProperly(t *testing.T) {
	container := component.NewContainer()
	ctx := context.Background()

	conditionContext := NewContext(ctx, container)
	assert.Equal(t, container, conditionContext.container)
	assert.Equal(t, ctx, conditionContext.ctx)
}

func TestNewConditionContextShouldPanicIfContextIsNotProvided(t *testing.T) {
	container := component.NewContainer()
	assert.PanicsWithValue(t, "nil context", func() {
		NewContext(nil, container)
	})
}

func TestNewConditionContextShouldPanicIfContainerIsNotProvided(t *testing.T) {
	assert.PanicsWithValue(t, "nil container", func() {
		NewContext(context.Background(), nil)
	})
}

func TestConditionContext_DeadlineShouldReturnWhenContextIsTimeout(t *testing.T) {
	container := component.NewContainer()
	now := time.Now()
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)

	conditionContext := NewContext(ctx, container)
	deadline, timeoutDefined := conditionContext.Deadline()
	assert.True(t, timeoutDefined)
	assert.Equal(t, now.Add(time.Second*3).Format(time.RFC3339), deadline.Format(time.RFC3339))
}

func TestConditionContext_DoneShouldWaitForContextToBeCompleted(t *testing.T) {
	container := component.NewContainer()
	now := time.Now()
	ctx, _ := context.WithTimeout(context.Background(), time.Second*1)

	conditionContext := NewContext(ctx, container)
	<-conditionContext.Done()
	assert.Equal(t, time.Now().Sub(now).Round(time.Second*1), time.Second*1)
}

func TestConditionContext_ErrShouldNotReturnErrorIfContextIsNotCancelled(t *testing.T) {
	container := component.NewContainer()

	conditionContext := NewContext(context.Background(), container)
	assert.Nil(t, conditionContext.Err())
}

func TestConditionContext_ErrShouldReturnErrorIfContextIsCancelled(t *testing.T) {
	container := component.NewContainer()
	ctx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()

	conditionContext := NewContext(ctx, container)
	assert.EqualError(t, conditionContext.Err(), "context canceled")
}

func TestConditionContext_ValueShouldReturnAssociatedValueWithKey(t *testing.T) {
	container := component.NewContainer()
	ctx := context.WithValue(context.Background(), "anyKey", "anyValue")

	conditionContext := NewContext(ctx, container)
	assert.Equal(t, "anyValue", conditionContext.Value("anyKey"))
}

func TestConditionContext_ContainerShouldReturnAnyContainerObject(t *testing.T) {
	container := component.NewContainer()

	conditionContext := NewContext(context.Background(), container)
	assert.Equal(t, container, conditionContext.Container())
}

func TestNewConditionEvaluatorShouldCreateEvaluatorProperly(t *testing.T) {
	container := component.NewContainer()

	evaluator := NewEvaluator(container)
	assert.Equal(t, container, evaluator.container)
}

func TestConditionEvaluator_EvaluateShouldReturnTrueIfAnyConditionIsNotProvided(t *testing.T) {
	container := component.NewContainer()

	evaluator := NewEvaluator(container)
	assert.True(t, evaluator.Evaluate(context.Background(), nil))
}

func TestConditionEvaluator_EvaluateShouldReturnTrueIfConditionMatch(t *testing.T) {
	container := component.NewContainer()
	ctx := context.Background()
	conditionContext := NewContext(ctx, container)

	mockCondition := &MockCondition{}
	mockCondition.On("MatchesCondition", conditionContext).Return(true)

	evaluator := NewEvaluator(container)
	assert.True(t, evaluator.Evaluate(ctx, []Condition{mockCondition}))
}

func TestConditionEvaluator_EvaluateShouldReturnFalseIfConditionDoesNotMatch(t *testing.T) {
	container := component.NewContainer()
	ctx := context.Background()
	conditionContext := NewContext(ctx, container)

	mockCondition := &MockCondition{}
	mockCondition.On("MatchesCondition", conditionContext).Return(false)

	evaluator := NewEvaluator(container)
	assert.False(t, evaluator.Evaluate(ctx, []Condition{mockCondition}))
}
