package component

import "context"

type ObjectProvider func(ctx context.Context) (any, error)
