package runtime

import (
	"codnect.io/procyon-core/runtime/property"
	"context"
)

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Port() int
}

type ServerProperties struct {
	property.Properties `prefix:"procyon.server"`

	Port int `prop:"port"`
}

func NewServerProperties() *ServerProperties {
	return &ServerProperties{}
}
