package runtime

import (
	"context"
	"errors"
	"github.com/codnect/procyoncore/event"
	"sync"
)

type EventMulticaster struct {
	listeners []event.Listener
	mu        sync.RWMutex
}

func NewEventMulticaster() *EventMulticaster {
	return &EventMulticaster{
		listeners: make([]event.Listener, 0),
	}
}

func (m *EventMulticaster) AddListener(eventListener event.Listener) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	if eventListener == nil {
		return errors.New("event listener cannot be nil")
	}

	m.listeners = append(m.listeners, eventListener)
	return nil
}

func (m *EventMulticaster) MulticastEvent(ctx context.Context, event event.Event) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	if ctx == nil {
		return errors.New("context cannot be nil")
	}

	if event == nil {
		return errors.New("event cannot be nil")
	}

	for _, listener := range m.listeners {
		if listener.Supports(event) {
			if listener.SupportsAsyncExecution() {
				// invoke listener async
			} else {
				err := listener.OnEvent(ctx, event)

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (m *EventMulticaster) RemoveListener(eventListener event.Listener) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	if eventListener == nil {
		return errors.New("event listener cannot be nil")
	}

	for index, listener := range m.listeners {
		if listener == eventListener {
			m.listeners = append(m.listeners[:index], m.listeners[index+1:]...)
			return nil
		}
	}

	return nil
}

func (m *EventMulticaster) RemoveAllListeners() {
	defer m.mu.Unlock()
	m.mu.Lock()
	m.listeners = make([]event.Listener, 0)
}
