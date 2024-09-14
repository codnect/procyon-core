package condition

import (
	"codnect.io/procyon-core/component"
	"codnect.io/procyon-core/component/filter"
	"codnect.io/procyon-core/runtime"
)

type OnPropertyCondition struct {
	name           string
	value          any
	matchIfMissing bool
}

func OnProperty(name string) *OnPropertyCondition {
	return &OnPropertyCondition{
		name: name,
	}
}

func (c *OnPropertyCondition) HavingValue(value any) *OnPropertyCondition {
	c.value = value
	return c
}

func (c *OnPropertyCondition) MatchIfMissing(matchIfMissing bool) *OnPropertyCondition {
	c.matchIfMissing = matchIfMissing
	return c
}

func (c *OnPropertyCondition) MatchesCondition(ctx component.ConditionContext) bool {
	container := ctx.Container()
	if container == nil {
		return false
	}

	result, err := container.GetObject(ctx, filter.ByTypeOf[runtime.Environment]())
	if err != nil {
		return false
	}

	environment := result.(runtime.Environment)
	value, ok := environment.PropertyResolver().Property(c.name)

	if !ok {
		return c.matchIfMissing
	}

	if c.value == nil {
		return true
	}

	return value == c.value
}
