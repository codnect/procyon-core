package condition

import (
	"codnect.io/procyon-core/container"
	"codnect.io/procyon-core/runtime/env"
)

type Evaluator interface {
	ShouldSkip(conditions []Condition) bool
}

type evaluator struct {
	ctx Context
}

func NewEvaluator(container container.Container, environment env.Environment) Evaluator {
	if container == nil {
		panic("condition: container cannot be nil")
	}

	if environment == nil {
		panic("condition: environment cannot be nil")
	}

	return &evaluator{
		newContext(container, environment),
	}
}

func (e *evaluator) ShouldSkip(conditions []Condition) bool {
	if len(conditions) == 0 {
		return false
	}

	for _, condition := range conditions {
		if !condition.Matches(e.ctx) {
			return false
		}
	}

	return true
}
