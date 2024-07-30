package condition

import "github.com/codnect/procyoncore/component"

type Evaluator struct {
	ctx Context
}

func NewEvaluator(definitionRegistry component.DefinitionRegistry) Evaluator {
	if definitionRegistry == nil {
		panic("definition registry cannot be nil")
	}

	return Evaluator{
		ctx: Context{
			definitionRegistry: definitionRegistry,
		},
	}
}

func (e Evaluator) ShouldSkip(conditions []Condition) bool {
	if len(conditions) == 0 {
		return false
	}

	for _, condition := range conditions {
		if !condition.Matches(e.ctx) {
			return false
		}
	}

	return true
}
