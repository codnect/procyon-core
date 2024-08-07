package component

import (
	"context"
	"time"
)

type CreatedEvent struct {
	object     any
	definition *Definition
	time       time.Time
}

type CreatedEventListener interface {
	OnCreated(ctx context.Context, event CreatedEvent) (any, error)
}

func newCreatedEvent(definition *Definition, object any) CreatedEvent {
	return CreatedEvent{
		object:     object,
		definition: definition,
		time:       time.Now(),
	}
}

func (e CreatedEvent) EventSource() any {
	return nil
}

func (e CreatedEvent) Time() time.Time {
	return e.time
}

func (e CreatedEvent) Object() any {
	return e.object
}

type InitializedEvent struct {
	object     any
	definition *Definition
	time       time.Time
}

type InitializedEventListener interface {
	OnInitialized(ctx context.Context, event InitializedEvent) (any, error)
}

func newInitializedEvent(definition *Definition, object any) InitializedEvent {
	return InitializedEvent{
		object:     object,
		definition: definition,
		time:       time.Now(),
	}
}

func (e InitializedEvent) EventSource() any {
	return nil
}

func (e InitializedEvent) Time() time.Time {
	return e.time
}

func (e InitializedEvent) Object() any {
	return e.object
}
