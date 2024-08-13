package runtime

import (
	"context"
	"github.com/codnect/procyoncore/component"
)

type Context interface {
	context.Context

	Environment() Environment
	Container() component.Container
}

type ContextConfigurer interface {
	ConfigureContext(ctx Context) error
}
