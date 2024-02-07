package runtime

import (
	"codnect.io/procyon-core/container"
	"codnect.io/procyon-core/runtime/env"
	"codnect.io/procyon-core/runtime/event"
	"context"
)

type Context interface {
	context.Context

	event.Publisher
	event.ListenerRegistry

	Environment() env.Environment
	Container() container.Container
	Close()
}

type ContextCustomizer interface {
	CustomizeContext(ctx Context) error
}
