package gameplatformerphase

import (
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/tilemaplayer"
)

// mockBodiesSpace implements body.BodiesSpace
type mockBodiesSpace struct {
	bodies []body.Collidable
}

func (m *mockBodiesSpace) AddBody(b body.Collidable) { m.bodies = append(m.bodies, b) }
func (m *mockBodiesSpace) Bodies() []body.Collidable { return m.bodies }
func (m *mockBodiesSpace) RemoveBody(b body.Collidable) {
	for i, c := range m.bodies {
		if c == b {
			m.bodies = append(m.bodies[:i], m.bodies[i+1:]...)
			return
		}
	}
}
func (m *mockBodiesSpace) QueueForRemoval(body.Collidable)                                     {}
func (m *mockBodiesSpace) ProcessRemovals()                                                    {}
func (m *mockBodiesSpace) Clear()                                                              {}
func (m *mockBodiesSpace) ResolveCollisions(body.Collidable) (bool, bool)                      { return false, false }
func (m *mockBodiesSpace) SetTilemapDimensionsProvider(tilemaplayer.TilemapDimensionsProvider) {}
func (m *mockBodiesSpace) GetTilemapDimensionsProvider() tilemaplayer.TilemapDimensionsProvider {
	return nil
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
