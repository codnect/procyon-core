package component

import (
	"context"
	"fmt"
	"reflect"
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
		typ := reflect.TypeOf(object)
		typeName := fmt.Sprintf("%s.%s", typ.PkgPath(), typ.Name())
		log.I(ctx, "Component '{}' is not eligible for ObjectProcessor", typeName)
	}

	return object, nil
}
