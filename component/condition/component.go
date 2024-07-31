package condition

import "github.com/codnect/procyoncore/component"

type OnMissingCondition struct {
}

func OnMissing(name string) *OnMissingCondition {
	return nil
}

func (c *OnMissingCondition) Matches(ctx component.ConditionContext) bool {
	return false
}
