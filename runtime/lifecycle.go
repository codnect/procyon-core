package runtime

import (
	"codnect.io/procyon-core/runtime/env/property"
	"context"
	"time"
)

type Lifecycle interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	IsRunning() bool
}

type LifecycleProperties struct {
	property.Properties `prefix:"procyon.lifecycle"`

	ShutdownTimeout time.Duration `prop:"shutdown.timeout" default:"30000"`
}

/*
type LifecycleManager struct {
	properties LifecycleProperties
	container  container.Container
}

func NewLifecycleManager(properties LifecycleProperties, container container.Container) *LifecycleManager {
	return &LifecycleManager{
		properties: properties,
		container:  container,
	}
}

func (p *LifecycleManager) start(ctx context.Context) error {
	err := p.startLifecycleComponents(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (p *LifecycleManager) stop(ctx context.Context) error {
	err := p.stopLifecycleComponents(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (p *LifecycleManager) startLifecycleComponents(ctx context.Context) (err error) {
	sharedInstances := p.container.SharedInstances()
	lifecycleInstances := sharedInstances.FindAllByType(reflector.TypeOf[Lifecycle]())

	for _, instance := range lifecycleInstances {
		lifecycle := instance.(Lifecycle)

		err = lifecycle.Start(ctx)

		if err != nil {
			return
		}
	}

	return
}

func (p *LifecycleManager) stopLifecycleComponents(ctx context.Context) (err error) {
	sharedInstances := p.container.SharedInstances()
	lifecycleInstances := sharedInstances.FindAllByType(reflector.TypeOf[Lifecycle]())

	for _, instance := range lifecycleInstances {
		lifecycle := instance.(Lifecycle)

		err = lifecycle.Stop(ctx)

		if err != nil {
			return
		}
	}

	return
}
*/
