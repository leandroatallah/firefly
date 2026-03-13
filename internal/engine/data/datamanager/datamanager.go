package datamanager

import "log"

// Manager holds the raw data for assets like JSON files.
type Manager struct {
	data map[string][]byte
}

// NewDataManager creates a new data manager.
func NewDataManager() *Manager {
	return &Manager{
		data: make(map[string][]byte),
	}
}

// Add stores the data for a given asset name.
func (m *Manager) Add(name string, data []byte) {
	m.data[name] = data
}

// Get retrieves the data for a given asset name.
func (m *Manager) Get(name string) []byte {
	data, ok := m.data[name]
	if !ok {
		log.Printf("data not found: %s", name)
		return nil
	}
	return data
}
