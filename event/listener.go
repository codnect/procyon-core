package event

import "context"

type Listener interface {
	OnEvent(ctx context.Context, event Event)
}

func Listen[E Event](fn func(ctx context.Context, event E) error) Listener {
	return nil
}
