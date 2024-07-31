package condition

import "github.com/codnect/procyoncore/component"

type OnTypeCondition struct {
}

func OnType[T any]() *OnTypeCondition {
	return nil
}

func (c *OnTypeCondition) Matches(ctx component.ConditionContext) bool {
	return false
}

type OnMissingTypeCondition struct {
}

func OnMissingType[T any]() *OnMissingTypeCondition {
	return nil
}

func (c *OnMissingTypeCondition) Matches(ctx component.ConditionContext) bool {
	return false
}
