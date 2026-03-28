package actors

import "fmt"

// Manager holds a registry of all active actors in a scene.
type Manager struct {
	actors        map[string]ActorEntity
	primaryActor  ActorEntity
}

// NewManager creates a new actor manager.
func NewManager() *Manager {
	return &Manager{
		actors: make(map[string]ActorEntity),
	}
}

// RegisterPrimary designates an actor as the primary (player-controlled) entity.
func (m *Manager) RegisterPrimary(actor ActorEntity) {
	m.primaryActor = actor
}

// Register adds an actor to the manager.
func (m *Manager) Register(actor ActorEntity) {
	id := actor.ID()
	if _, exists := m.actors[id]; exists {
		fmt.Printf("Warning: Actor with ID '%s' is already registered. Overwriting.\n", id)
	}
	m.actors[id] = actor
}

// Find retrieves an actor by its ID.
func (m *Manager) Find(id string) (ActorEntity, bool) {
	actor, found := m.actors[id]
	return actor, found
}

// Unregister removes an actor from the manager.
func (m *Manager) Unregister(actor ActorEntity) {
	delete(m.actors, actor.ID())
}

// Clear removes all registered actors.
func (m *Manager) Clear() {
	for k := range m.actors {
		delete(m.actors, k)
	}
}

// GetPlayer retrieves the primary actor registered via RegisterPrimary.
// Falls back to finding an actor with ID "player" for backward compatibility.
func (m *Manager) GetPlayer() (ActorEntity, bool) {
	if m.primaryActor != nil {
		return m.primaryActor, true
	}
	return m.Find("player")
}

func (m *Manager) ForEach(f func(ActorEntity)) {
	for _, actor := range m.actors {
		f(actor)
	}
}
