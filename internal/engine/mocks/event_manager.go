package mocks

// MockEventManager is a shared mock for the EventManager interface.
type MockEventManager struct {
	PublishFunc func(e interface{})
}

func (m *MockEventManager) Publish(e interface{}) {
	if m.PublishFunc != nil {
		m.PublishFunc(e)
	}
}
