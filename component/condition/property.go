package condition

import "github.com/codnect/procyoncore/component"

type OnPropertyCondition struct {
}

func OnProperty(name string) *OnPropertyCondition {
	return nil
}

func (c *OnPropertyCondition) Prefix(prefix string) *OnPropertyCondition {
	return c
}

func (c *OnPropertyCondition) HavingValue(value string) *OnPropertyCondition {
	return c
}

func (c *OnPropertyCondition) MatchIfMissing(matchIfMissing bool) *OnPropertyCondition {
	return c
}

func (c *OnPropertyCondition) Matches(ctx component.ConditionContext) bool {
	return false
}
