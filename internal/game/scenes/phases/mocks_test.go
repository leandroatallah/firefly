package gamescenephases

import (
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/audio"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/navigation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/sequences"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/tilemaplayer"
	"github.com/hajimehoshi/ebiten/v2"
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

// mockSequencePlayer implements sequences.Player
type mockSequencePlayer struct {
	playing bool
}

func (m *mockSequencePlayer) IsPlaying() bool           { return m.playing }
func (m *mockSequencePlayer) IsOver() bool              { return !m.playing }
func (m *mockSequencePlayer) Play(_ sequences.Sequence) {}
func (m *mockSequencePlayer) PlaySequence(_ string)     {}
func (m *mockSequencePlayer) Stop()                     {}
func (m *mockSequencePlayer) Update()                   {}

// mockGoal implements phases.Goal
type mockGoal struct {
	completed          bool
	onCompletionCalled bool
}

func (m *mockGoal) IsCompleted() bool { return m.completed }
func (m *mockGoal) OnCompletion()     { m.onCompletionCalled = true }

// mockSceneManager implements navigation.SceneManager
type mockSceneManager struct {
	navigateToCalled   bool
	navigateBackCalled bool
}

func (m *mockSceneManager) AudioManager() *audio.AudioManager { return nil }
func (m *mockSceneManager) Draw(_ *ebiten.Image)              {}
func (m *mockSceneManager) NavigateTo(_ navigation.SceneType, _ navigation.Transition, _ bool) {
	m.navigateToCalled = true
}
func (m *mockSceneManager) NavigateBack(_ navigation.Transition) { m.navigateBackCalled = true }
func (m *mockSceneManager) SwitchTo(_ navigation.Scene)          {}
func (m *mockSceneManager) Update() error                        { return nil }
func (m *mockSceneManager) CurrentScene() navigation.Scene       { return nil }
func (m *mockSceneManager) IsTransitioning() bool                { return false }

// mockCollidable implements body.Collidable (minimal stub)
type mockCollidable struct {
	id string
}

func (c *mockCollidable) ID() string                                          { return c.id }
func (c *mockCollidable) SetID(id string)                                     { c.id = id }
func (c *mockCollidable) Owner() interface{}                                  { return nil }
func (c *mockCollidable) SetOwner(interface{})                                {}
func (c *mockCollidable) LastOwner() interface{}                              { return nil }
func (c *mockCollidable) Position() image.Rectangle                           { return image.Rectangle{} }
func (c *mockCollidable) SetPosition(_, _ int)                                {}
func (c *mockCollidable) SetPosition16(_, _ int)                              {}
func (c *mockCollidable) SetSize(_, _ int)                                    {}
func (c *mockCollidable) Scale() float64                                      { return 1 }
func (c *mockCollidable) SetScale(float64)                                    {}
func (c *mockCollidable) GetPosition16() (int, int)                           { return 0, 0 }
func (c *mockCollidable) GetPositionMin() (int, int)                          { return 0, 0 }
func (c *mockCollidable) GetShape() body.Shape                                { return nil }
func (c *mockCollidable) GetTouchable() body.Touchable                        { return nil }
func (c *mockCollidable) DrawCollisionBox(_ *ebiten.Image, _ image.Rectangle) {}
func (c *mockCollidable) CollisionPosition() []image.Rectangle                { return nil }
func (c *mockCollidable) CollisionShapes() []body.Collidable                  { return nil }
func (c *mockCollidable) IsObstructive() bool                                 { return false }
func (c *mockCollidable) SetIsObstructive(bool)                               {}
func (c *mockCollidable) AddCollision(...body.Collidable)                     {}
func (c *mockCollidable) ClearCollisions()                                    {}
func (c *mockCollidable) SetTouchable(body.Touchable)                         {}
func (c *mockCollidable) ApplyValidPosition(_ int, _ bool, _ body.BodiesSpace) (int, int, bool) {
	return 0, 0, false
}
func (c *mockCollidable) OnTouch(body.Collidable) {}
func (c *mockCollidable) OnBlock(body.Collidable) {}
