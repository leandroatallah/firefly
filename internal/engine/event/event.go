package event

// Event defines the interface for all events.
type Event interface {
	Type() string
}

// Listener is a function that handles an event.
type Listener func(e Event)

// wrappedListener holds a listener with its unique ID.
type wrappedListener struct {
	id       int
	listener Listener
}

// Manager handles event subscription and dispatching.
type Manager struct {
	listeners map[string][]wrappedListener
	nextID    int
}

// NewManager creates a new event manager.
func NewManager() *Manager {
	return &Manager{
		listeners: make(map[string][]wrappedListener),
		nextID:    1,
	}
}

// Subscribe adds a listener for a given event type and returns an unsubscribe function.
// The returned function removes the listener when called.
func (m *Manager) Subscribe(eventType string, listener Listener) func() {
	id := m.nextID
	m.nextID++
	m.listeners[eventType] = append(m.listeners[eventType], wrappedListener{id: id, listener: listener})

	return func() {
		listeners := m.listeners[eventType]
		for i, wl := range listeners {
			if wl.id == id {
				m.listeners[eventType] = append(listeners[:i], listeners[i+1:]...)
				return
			}
		}
	}
}

// Publish dispatches an event to all registered listeners.
func (m *Manager) Publish(e Event) {
	if listeners, ok := m.listeners[e.Type()]; ok {
		for _, wl := range listeners {
			wl.listener(e)
		}
	}
}

// GenericEvent is a simple event implementation that holds a type and a payload.
type GenericEvent struct {
	EventType string
	Payload   map[string]interface{}
}

// Type returns the event type.
func (e GenericEvent) Type() string {
	return e.EventType
}
