package runtime

import (
	"codnect.io/procyon-core/component"
	"context"
)

type Context interface {
	context.Context

	Start() error
	Stop() error
	IsRunning() bool

	Environment() Environment
	Container() component.Container
}

type ContextConfigurer interface {
	ConfigureContext(ctx Context) error
}
