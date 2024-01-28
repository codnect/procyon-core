package runtime

import (
	"codnect.io/procyon-core/container"
	"codnect.io/procyon-core/event"
	"codnect.io/procyon-core/runtime/env"
	"context"
	"time"
)

type Context interface {
	context.Context
	event.Publisher
	event.ListenerRegistry

	ApplicationName() string
	DisplayName() string
	StartupTime() time.Time
	Environment() env.Environment
	Container() container.Container

	Start() error
	Stop() error
}

type ContextCustomizer interface {
	CustomizeContext(ctx Context) error
}
