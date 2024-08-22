package event

/*
type Multicaster interface {
	ListenerRegistry

	MulticastEvent(ctx context.Context, event ApplicationEvent) error
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

func (m *SimpleMulticaster) AddListener(eventListener Listener) error {
	defer m.mu.Unlock()
	m.mu.Lock()

	if eventListener == nil {
		return errors.New("event listener cannot be nil")
	}

	m.listeners = append(m.listeners, eventListener)
	return nil
}

func (m *SimpleMulticaster) MulticastEvent(ctx context.Context, event ApplicationEvent) error {
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

func (m *SimpleMulticaster) RemoveListener(eventListener Listener) error {
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

func (m *SimpleMulticaster) RemoveAllListeners() {
	defer m.mu.Unlock()
	m.mu.Lock()
	m.listeners = make([]Listener, 0)
}

*/
