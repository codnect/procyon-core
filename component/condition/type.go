package condition

import (
	"codnect.io/procyon-core/component"
	"codnect.io/procyon-core/component/filter"
	"reflect"
)

type OnTypeCondition struct {
	typ reflect.Type
}

func OnType[T any]() *OnTypeCondition {
	return &OnTypeCondition{
		typ: reflect.TypeFor[T](),
	}
}

func (c *OnTypeCondition) MatchesCondition(ctx component.ConditionContext) bool {
	container := ctx.Container()
	definitions := container.Definitions().List(filter.ByType(c.typ))
	singletons := container.Singletons().List(filter.ByType(c.typ))
	return len(definitions) != 0 || len(singletons) != 0
}

type OnMissingTypeCondition struct {
	missingType reflect.Type
}

func OnMissingType[T any]() *OnMissingTypeCondition {
	return &OnMissingTypeCondition{
		missingType: reflect.TypeFor[T](),
	}
}

func (c *OnMissingTypeCondition) MatchesCondition(ctx component.ConditionContext) bool {
	container := ctx.Container()
	definitions := container.Definitions().List(filter.ByType(c.missingType))
	singletons := container.Singletons().List(filter.ByType(c.missingType))
	return len(definitions) == 0 && len(singletons) == 0
}
