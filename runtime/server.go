package runtime

import "context"

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Port() int
}
