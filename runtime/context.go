package runtime

import (
	"codnect.io/procyon-core/component"
	"context"
)

type Context interface {
	context.Context

	Environment() Environment
	Container() component.Container
}

type ContextConfigurer interface {
	ConfigureContext(ctx Context) error
}
