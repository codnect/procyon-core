package event

import "context"

type TypedPublisher[E ApplicationEvent] interface {
	PublishEvent(ctx context.Context, event E) error
}

type AsyncTypedPublisher[E ApplicationEvent] interface {
	PublishEventAsync(ctx context.Context, event E) error
}

type Publisher interface {
	PublishEvent(ctx context.Context, event ApplicationEvent) error
	PublishEventAsync(ctx context.Context, event ApplicationEvent) error
}
