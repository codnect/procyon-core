package event

import "context"

type Multicaster interface {
	AddListener(listener Listener) error
	MulticastEvent(ctx context.Context, event Event) error
	RemoveListener(listener Listener) error
	RemoveAllListeners()
}
