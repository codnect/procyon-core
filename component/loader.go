package component

import "context"

type DefinitionLoader struct {
	container Container
	evaluator ConditionEvaluator
}

func NewDefinitionLoader(container Container) *DefinitionLoader {
	return &DefinitionLoader{
		container: container,
		evaluator: NewConditionEvaluator(container),
	}
}

func (l *DefinitionLoader) LoadDefinitions(ctx context.Context, components []*Component) error {
	skippedComponents := make([]*Component, 0)

	for _, component := range components {
		if l.evaluator.Evaluate(ctx, component.Conditions()) {
			err := l.container.Definitions().Register(component.Definition())

			if err != nil {
				return err
			}
		} else {
			skippedComponents = append(skippedComponents, component)
		}
	}

	if len(components) == len(skippedComponents) {
		return nil
	}

	return l.LoadDefinitions(ctx, skippedComponents)
}
