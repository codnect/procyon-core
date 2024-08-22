package event

import (
	"context"
	"errors"
)

type Listener interface {
	OnEvent(ctx context.Context, event ApplicationEvent) error
	SupportsEvent(event ApplicationEvent) bool
}

func Listen[E ApplicationEvent](handler func(ctx context.Context, event E) error) Listener {
	return listenerAdapter[E]{
		handler: handler,
	}
}

type listenerAdapter[E any] struct {
	handler func(ctx context.Context, event E) error
}

func (a listenerAdapter[E]) OnEvent(ctx context.Context, event ApplicationEvent) error {
	if !a.SupportsEvent(event) {
		return errors.New("")
	}

	return a.handler(ctx, event.(E))
}

func (a listenerAdapter[E]) SupportsEvent(event ApplicationEvent) bool {
	if _, ok := event.(E); ok {
		return true
	}

	return false
}
