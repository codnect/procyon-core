package event

import (
	"context"
	"errors"
)

type Listener interface {
	OnEvent(ctx context.Context, event Event) error
	Supports(event Event) bool
	SupportsAsyncExecution() bool
}

func Listen[E Event](handler func(ctx context.Context, event E) error) Listener {
	return handlerWrapper[E]{
		handler: handler,
		isAsync: false,
	}
}

func ListenAsync[E Event](handler func(ctx context.Context, event E) error) Listener {
	return handlerWrapper[E]{
		handler: handler,
		isAsync: true,
	}
}

type handlerWrapper[E any] struct {
	handler func(ctx context.Context, event E) error
	isAsync bool
}

func (w handlerWrapper[E]) OnEvent(ctx context.Context, event Event) error {
	if !w.Supports(event) {
		return errors.New("")
	}

	return w.handler(ctx, event.(E))
}

func (w handlerWrapper[E]) Supports(event Event) bool {
	if _, ok := event.(E); ok {
		return true
	}

	return false
}

func (w handlerWrapper[E]) SupportsAsyncExecution() bool {
	return w.isAsync
}
