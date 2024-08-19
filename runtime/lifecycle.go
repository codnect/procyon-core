package runtime

import (
	"codnect.io/procyon-core/runtime/property"
	"context"
	"time"
)

type Lifecycle interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	IsRunning() bool
}

type LifecycleProperties struct {
	property.Properties `prefix:"procyon.lifecycle"`

	ShutdownTimeout time.Duration `prop:"shutdown-timeout" default:"30000"`
}

func NewLifecycleProperties() *ServerProperties {
	return &ServerProperties{}
}
