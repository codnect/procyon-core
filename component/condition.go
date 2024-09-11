package component

import (
	"context"
	"time"
)

type Condition interface {
	MatchesCondition(ctx ConditionContext) bool
}

type ConditionContext struct {
	ctx       context.Context
	container Container
}

func NewConditionContext(ctx context.Context, container Container) ConditionContext {
	if ctx == nil {
		panic("nil context")
	}

	if container == nil {
		panic("nil container")
	}

	return ConditionContext{
		ctx:       ctx,
		container: container,
	}
}

func (c ConditionContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c ConditionContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c ConditionContext) Err() error {
	return c.ctx.Err()
}

func (c ConditionContext) Value(key any) any {
	return c.ctx.Value(key)
}

func (c ConditionContext) Container() Container {
	return c.container
}

type ConditionEvaluator struct {
	container Container
}

func NewConditionEvaluator(container Container) ConditionEvaluator {
	if container == nil {
		panic("nil container")
	}

	return ConditionEvaluator{
		container: container,
	}
}

func (e ConditionEvaluator) Evaluate(ctx context.Context, conditions []Condition) bool {
	if len(conditions) == 0 {
		return true
	}

	conditionContext := NewConditionContext(ctx, e.container)

	for _, condition := range conditions {
		if !condition.MatchesCondition(conditionContext) {
			return false
		}
	}

	return true
}
