package component

import (
	"context"
	"time"
)

type CreationContext struct {
	ctx        context.Context
	definition *Definition
	object     any
}

func NewCreationContext(ctx context.Context, definition *Definition, object any) CreationContext {
	return CreationContext{
		ctx: ctx,
	}
}

func (c CreationContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c CreationContext) Done() <-chan struct{} {
	return c.Done()
}

func (c CreationContext) Err() error {
	return c.Err()
}

func (c CreationContext) Value(key any) any {
	return c.Value(key)
}

func (c CreationContext) Definition() *Definition {
	return c.definition
}

func (c CreationContext) Object() any {
	return c.object
}

type BeforeInitialization interface {
	BeforeInit(ctx CreationContext) (any, error)
}

type AfterInitialization interface {
	AfterInit(ctx CreationContext) (any, error)
}

type Initialization interface {
	DoInit(ctx context.Context) error
}
