package event

import (
	"context"
	"errors"
)

type Listener interface {
	OnEvent(ctx context.Context, event ApplicationEvent) error
	Supports(event ApplicationEvent) bool
}

type ListenerRegistry interface {
	AddListener(listener Listener) error
	RemoveListener(listener Listener) error
	RemoveAllListeners()
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
	if !a.Supports(event) {
		return errors.New("")
	}

	return a.handler(ctx, event.(E))
}

func (a listenerAdapter[E]) Supports(event ApplicationEvent) bool {
	if _, ok := event.(E); ok {
		return true
	}

	return false
}
