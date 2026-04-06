package projectile

import (
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/tilemaplayer"
)

// mockBodiesSpace implements body.BodiesSpace for testing.
type mockBodiesSpace struct {
	bodies              []body.Collidable
	AddBodyFunc         func(body.Collidable)
	QueueForRemovalFunc func(body.Collidable)
	RemoveBodyFunc      func(body.Collidable)
	tilemapProvider     tilemaplayer.TilemapDimensionsProvider
	queuedForRemoval    []body.Collidable
}

func (m *mockBodiesSpace) AddBody(b body.Collidable) {
	if m.AddBodyFunc != nil {
		m.AddBodyFunc(b)
	}
	m.bodies = append(m.bodies, b)
}
func (m *mockBodiesSpace) Bodies() []body.Collidable { return m.bodies }
func (m *mockBodiesSpace) RemoveBody(b body.Collidable) {
	if m.RemoveBodyFunc != nil {
		m.RemoveBodyFunc(b)
	}
	for i, c := range m.bodies {
		if c == b {
			m.bodies = append(m.bodies[:i], m.bodies[i+1:]...)
			return
		}
	}
}
func (m *mockBodiesSpace) QueueForRemoval(b body.Collidable) {
	if m.QueueForRemovalFunc != nil {
		m.QueueForRemovalFunc(b)
	}
	m.queuedForRemoval = append(m.queuedForRemoval, b)
}
func (m *mockBodiesSpace) ProcessRemovals() {
	for _, b := range m.queuedForRemoval {
		m.RemoveBody(b)
	}
	m.queuedForRemoval = nil
}
func (m *mockBodiesSpace) Clear()                                         {}
func (m *mockBodiesSpace) ResolveCollisions(body.Collidable) (bool, bool) { return false, false }
func (m *mockBodiesSpace) SetTilemapDimensionsProvider(p tilemaplayer.TilemapDimensionsProvider) {
	m.tilemapProvider = p
}
func (m *mockBodiesSpace) GetTilemapDimensionsProvider() tilemaplayer.TilemapDimensionsProvider {
	return m.tilemapProvider
}
func (m *mockBodiesSpace) Find(id string) body.Collidable {
	for _, b := range m.bodies {
		if b.ID() == id {
			return b
		}
	}
	return nil
}
func (m *mockBodiesSpace) Query(image.Rectangle) []body.Collidable { return nil }

// mockTilemapDimensionsProvider implements tilemaplayer.TilemapDimensionsProvider.
type mockTilemapDimensionsProvider struct {
	width, height int
}

func (p *mockTilemapDimensionsProvider) GetTilemapWidth() int  { return p.width }
func (p *mockTilemapDimensionsProvider) GetTilemapHeight() int { return p.height }
func (p *mockTilemapDimensionsProvider) GetCameraBounds() (image.Rectangle, bool) {
	return image.Rectangle{}, false
}
