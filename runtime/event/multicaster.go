package event

import (
	"context"
	"errors"
	"sync"
)

type Multicaster interface {
	AddEventListener(listener Listener) error
	MulticastEvent(ctx context.Context, event ApplicationEvent) error
	MulticastEventAsync(ctx context.Context, event ApplicationEvent) error
	RemoveEventListener(listener Listener) error
	RemoveAllEventListeners()
}

type SimpleMulticaster struct {
	listeners []Listener
	mu        sync.RWMutex
}

func NewSimpleMulticaster() *SimpleMulticaster {
	return &SimpleMulticaster{
		listeners: make([]Listener, 0),
	}
}

func (m *SimpleMulticaster) AddEventListener(listener Listener) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	if listener == nil {
		return errors.New("nil event listener")
	}

	m.listeners = append(m.listeners, listener)
	return nil
}

func (m *SimpleMulticaster) MulticastEvent(ctx context.Context, event ApplicationEvent) error {
	m.mu.Lock()
	listeners := make([]Listener, len(m.listeners))
	copy(listeners, m.listeners)
	m.mu.Unlock()

	if ctx == nil {
		return errors.New("nil context")
	}

	if event == nil {
		return errors.New("nil event")
	}

	for _, listener := range listeners {
		if listener.SupportsEvent(event) {
			err := listener.OnEvent(ctx, event)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *SimpleMulticaster) MulticastEventAsync(ctx context.Context, event ApplicationEvent) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	if ctx == nil {
		return errors.New("nil context")
	}

	if event == nil {
		return errors.New("nil event")
	}

	for _, listener := range m.listeners {
		if listener.SupportsEvent(event) {
			go func() {
				err := listener.OnEvent(ctx, event)
				if err != nil {
					// handle error
				}
			}()
		}
	}

	return nil
}

func (m *SimpleMulticaster) RemoveEventListener(eventListener Listener) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	if eventListener == nil {
		return errors.New("nil event listener")
	}

	for index, listener := range m.listeners {
		if listener == eventListener {
			m.listeners = append(m.listeners[:index], m.listeners[index+1:]...)
			return nil
		}
	}

	return nil
}

func (m *SimpleMulticaster) RemoveAllEventListeners() {
	defer m.mu.Unlock()
	m.mu.Lock()
	m.listeners = make([]Listener, 0)
}
