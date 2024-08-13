package condition

import (
	"codnect.io/procyon-core/component"
	"codnect.io/procyon-core/component/filter"
	"codnect.io/reflector"
)

type OnTypeCondition struct {
	typ reflector.Type
}

func OnType[T any]() *OnTypeCondition {
	return &OnTypeCondition{
		typ: reflector.TypeOf[T](),
	}
}

func (c *OnTypeCondition) MatchesCondition(ctx component.ConditionContext) bool {
	container := ctx.Container()
	definitions := container.Definitions().List(filter.ByType(c.typ))

	if len(definitions) != 0 {
		return true
	}

	return false
}

type OnMissingTypeCondition struct {
	missingType reflector.Type
}

func OnMissingType[T any]() *OnMissingTypeCondition {
	return &OnMissingTypeCondition{
		missingType: reflector.TypeOf[T](),
	}
}

func (c *OnMissingTypeCondition) MatchesCondition(ctx component.ConditionContext) bool {
	container := ctx.Container()
	definitions := container.Definitions().List(filter.ByType(c.missingType))

	if len(definitions) == 0 {
		return true
	}

	return false
}
