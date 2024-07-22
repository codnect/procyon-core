package component

type ConditionContext interface {
}

type Condition interface {
	Matches(ctx ConditionContext) bool
}

type ConditionEvaluator interface {
	ShouldSkip(conditions []Condition) bool
}
