package actors

import "fmt"

// Manager holds a registry of all active actors in a scene.
type Manager struct {
	actors map[string]ActorEntity
}

// NewManager creates a new actor manager.
func NewManager() *Manager {
	return &Manager{
		actors: make(map[string]ActorEntity),
	}
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

// GetPlayer retrieves the player actor.
// It assumes the player is registered with the ID "player".
func (m *Manager) GetPlayer() (PlayerEntity, bool) {
	actor, found := m.Find("player")
	if !found {
		return nil, false
	}
	player, ok := actor.(PlayerEntity)
	if !ok {
		// This should ideally not happen if "player" ID is reserved for PlayerEntity
		fmt.Printf("Warning: Actor with ID 'player' is not a PlayerEntity.\n")
		return nil, false
	}
	return player, true
}
