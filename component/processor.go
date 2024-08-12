package component

import (
	"codnect.io/reflector"
	"context"
)

type ObjectProcessor interface {
	ProcessBeforeInit(ctx context.Context, object any) (any, error)
	ProcessAfterInit(ctx context.Context, object any) (any, error)
}

type processorChecker struct {
	container      *ObjectContainer
	processorCount int
}

func newProcessorChecker(container *ObjectContainer, processorCount int) processorChecker {
	return processorChecker{
		container:      container,
		processorCount: processorCount,
	}
}

func (c processorChecker) ProcessBeforeInit(ctx context.Context, object any) (any, error) {
	return object, nil
}

func (c processorChecker) ProcessAfterInit(ctx context.Context, object any) (any, error) {
	if _, ok := object.(ObjectProcessor); !ok && c.container.ObjectProcessorCount() < c.processorCount {
		log.I(ctx, "Component '{}' is not eligible for ObjectProcessor", fullTypeName(reflector.TypeOfAny(object)))
	}

	return object, nil
}
