package component

type Condition interface {
	Matches(ctx ConditionContext) bool
}

type ConditionContext struct {
	definitionRegistry DefinitionRegistry
}

type ConditionEvaluator struct {
	ctx ConditionContext
}

func NewConditionEvaluator(definitionRegistry DefinitionRegistry) ConditionEvaluator {
	if definitionRegistry == nil {
		panic("definition registry cannot be nil")
	}

	return ConditionEvaluator{
		ctx: ConditionContext{
			definitionRegistry: definitionRegistry,
		},
	}
}

func (e ConditionEvaluator) Evaluate(conditions []Condition) bool {
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
