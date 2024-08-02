package runtime

import (
	"context"
	"github.com/codnect/procyoncore/component"
)

type Context interface {
	context.Context

	Environment() Environment
	Container() component.Container
	Close()
}

type ContextCustomizer interface {
	CustomizeContext(ctx Context) error
}
