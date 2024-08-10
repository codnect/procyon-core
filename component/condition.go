package component

import (
	"context"
	"time"
)

type Condition interface {
	Matches(ctx ConditionContext) bool
}

type ConditionContext struct {
	ctx      context.Context
	registry DefinitionRegistry
}

func newConditionContext(ctx context.Context, registry DefinitionRegistry) ConditionContext {
	return ConditionContext{
		ctx:      ctx,
		registry: registry,
	}
}

func (c ConditionContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c ConditionContext) Done() <-chan struct{} {
	return c.Done()
}

func (c ConditionContext) Err() error {
	return c.Err()
}

func (c ConditionContext) Value(key any) any {
	return c.Value(key)
}

func (c ConditionContext) DefinitionRegistry() DefinitionRegistry {
	return c.registry
}

type ConditionEvaluator struct {
	registry DefinitionRegistry
}

func NewConditionEvaluator(registry DefinitionRegistry) ConditionEvaluator {
	if registry == nil {
		panic("definition registry cannot be nil")
	}

	return ConditionEvaluator{
		registry: registry,
	}
}

func (e ConditionEvaluator) Evaluate(ctx context.Context, conditions []Condition) bool {
	if len(conditions) == 0 {
		return false
	}

	conditionContext := newConditionContext(ctx, e.registry)

	for _, condition := range conditions {
		if !condition.Matches(conditionContext) {
			return false
		}
	}

	return true
}
