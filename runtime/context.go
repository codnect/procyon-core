package runtime

import (
	"codnect.io/procyon-core/component"
	"codnect.io/procyon-core/runtime/event"
	"context"
)

type Context interface {
	context.Context
	event.Publisher

	Start() error
	Stop() error
	IsRunning() bool

	Environment() Environment
	Container() component.Container
}

type ContextConfigurer interface {
	ConfigureContext(ctx Context) error
}
