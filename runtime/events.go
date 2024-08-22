package runtime

import "time"

type StartupEvent struct {
	ctx  Context
	time time.Time
}

func NewStartupEvent(ctx Context) StartupEvent {
	return StartupEvent{
		ctx:  ctx,
		time: time.Now(),
	}
}

func (s StartupEvent) Context() Context {
	return s.ctx
}

func (s StartupEvent) EventSource() any {
	return s.ctx
}

func (s StartupEvent) EventTime() time.Time {
	return s.time
}

type ShutdownEvent struct {
	ctx  Context
	time time.Time
}

func NewShutdownEvent(ctx Context) ShutdownEvent {
	return ShutdownEvent{
		ctx:  ctx,
		time: time.Now(),
	}
}

func (s ShutdownEvent) Context() Context {
	return s.ctx
}

func (s ShutdownEvent) EventSource() any {
	return s.ctx
}

func (s ShutdownEvent) EventTime() time.Time {
	return s.time
}
