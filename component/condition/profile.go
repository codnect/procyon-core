package condition

import "github.com/codnect/procyoncore/component"

type OnProfileCondition struct {
}

func OnProfile(profiles ...string) *OnProfileCondition {
	return nil
}

func (c *OnProfileCondition) Matches(ctx component.ConditionContext) bool {
	return false
}
