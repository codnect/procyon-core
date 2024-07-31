package runtime

import (
	"context"
	"github.com/codnect/procyoncore/component"
	"github.com/codnect/procyoncore/runtime/env"
)

type Context interface {
	context.Context

	Environment() env.Environment
	Container() component.Container
	Close()
}

type ContextCustomizer interface {
	CustomizeContext(ctx Context) error
}
