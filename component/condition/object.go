package condition

import (
	"github.com/codnect/procyoncore/component"
)

type OnObjectCondition struct {
	name string
}

func OnObject(name string) *OnObjectCondition {
	return &OnObjectCondition{
		name: name,
	}
}

func (c *OnObjectCondition) MatchesCondition(ctx component.ConditionContext) bool {
	container := ctx.Container()
	return container.Definitions().Contains(c.name) || container.Singletons().Contains(c.name)
}

type OnMissingObjectCondition struct {
	name string
}

func OnMissingObject(name string) *OnMissingObjectCondition {
	return &OnMissingObjectCondition{
		name: name,
	}
}

func (c *OnMissingObjectCondition) MatchesCondition(ctx component.ConditionContext) bool {
	container := ctx.Container()
	return !container.Definitions().Contains(c.name) && !container.Singletons().Contains(c.name)
}
