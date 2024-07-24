package component

import "context"

const (
	SingletonScope = "singleton"
	PrototypeScope = "prototype"
)

type Scope interface {
	GetObject(ctx context.Context, name string, provider ObjectProvider) (any, error)
	RemoveObject(ctx context.Context, name string) (any, error)
}
