package runtime

import (
	"codnect.io/procyon-core/runtime/env"
	"time"
)

type StartupListener interface {
	OnStarting(ctx Context)
	OnEnvironmentPrepared(ctx Context, environment env.Environment)
	OnContextPrepared(ctx Context)
	OnContextLoaded(ctx Context)
	OnContextStarted(ctx Context)
	OnStarted(ctx Context, timeTaken time.Duration)
	OnReady(ctx Context, timeTaken time.Duration)
	OnFailed(ctx Context, err error)
}

type StartupListeners []StartupListener

func (l StartupListeners) Starting(ctx Context) {
	for _, listener := range l {
		listener.OnStarting(ctx)
	}
}

func (l StartupListeners) EnvironmentPrepared(ctx Context, environment env.Environment) {
	for _, listener := range l {
		listener.OnEnvironmentPrepared(ctx, environment)
	}
}

func (l StartupListeners) ContextPrepared(ctx Context) {
	for _, listener := range l {
		listener.OnContextPrepared(ctx)
	}
}

func (l StartupListeners) ContextLoaded(ctx Context) {
	for _, listener := range l {
		listener.OnContextLoaded(ctx)
	}
}

func (l StartupListeners) ContextStarted(ctx Context) {
	for _, listener := range l {
		listener.OnContextStarted(ctx)
	}
}

func (l StartupListeners) Started(ctx Context, timeTaken time.Duration) {
	for _, listener := range l {
		listener.OnStarted(ctx, timeTaken)
	}
}

func (l StartupListeners) Ready(ctx Context, timeTaken time.Duration) {
	for _, listener := range l {
		listener.OnReady(ctx, timeTaken)
	}
}

func (l StartupListeners) Failed(ctx Context, err error) {
	for _, listener := range l {
		listener.OnFailed(ctx, err)
	}
}
