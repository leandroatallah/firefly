package projectile

import (
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/assets/font"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	contractscombat "github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/tilemaplayer"
	"github.com/boilerplate/ebiten-template/internal/engine/render/camera"
	"github.com/boilerplate/ebiten-template/internal/engine/render/particles"
	"github.com/boilerplate/ebiten-template/internal/kit/combat"
	"github.com/hajimehoshi/ebiten/v2"
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

// spawnPuffCall records a single SpawnPuff invocation for assertions.
type spawnPuffCall struct {
	typeKey   string
	x, y      float64
	count     int
	randRange float64
}

// mockVFXManager implements vfx.Manager for testing.
type mockVFXManager struct {
	spawnPuffCallCount int
	lastTypeKey        string
	lastX, lastY       float64
	spawnPuffCalls     []spawnPuffCall
	lastCount          int
	lastRandRange      float64
}

func (m *mockVFXManager) SetAppContext(_ any)                                 {}
func (m *mockVFXManager) Update()                                             {}
func (m *mockVFXManager) Draw(_ *ebiten.Image, _ *camera.Controller)          {}
func (m *mockVFXManager) AddParticle(_ *particles.Particle)                   {}
func (m *mockVFXManager) AddTrauma(_ *camera.Controller, _ float64)           {}
func (m *mockVFXManager) PixelConfig() *particles.Config                      { return nil }
func (m *mockVFXManager) SetDefaultFont(_ *font.FontText)                     {}
func (m *mockVFXManager) SpawnDeathExplosion(_, _ float64, _ int)             {}
func (m *mockVFXManager) SpawnFallingRocks(_, _, _ float64, _ int)            {}
func (m *mockVFXManager) SpawnFloatingText(_ string, _, _ float64, _ int)     {}
func (m *mockVFXManager) SpawnFloatingTextAbove(_ body.Body, _ string, _ int) {}
func (m *mockVFXManager) SpawnJumpPuff(_, _ float64, _ int)                   {}
func (m *mockVFXManager) SpawnLandingPuff(_, _ float64, _ int)                {}
func (m *mockVFXManager) SpawnPuff(typeKey string, x, y float64, count int, randRange float64) {
	m.spawnPuffCallCount++
	m.lastTypeKey = typeKey
	m.lastX = x
	m.lastY = y
	m.lastCount = count
	m.lastRandRange = randRange
	m.spawnPuffCalls = append(m.spawnPuffCalls, spawnPuffCall{
		typeKey:   typeKey,
		x:         x,
		y:         y,
		count:     count,
		randRange: randRange,
	})
}
func (m *mockVFXManager) SpawnDirectionalPuff(_ string, _ float64, _ float64, _ bool, _ int, _ float64) {
}
func (m *mockVFXManager) TriggerScreenFlash() {}

// mockCollidable implements body.Collidable for testing.
type mockCollidable struct {
	id       string
	owner    interface{}
	x16, y16 int
}

func (m *mockCollidable) ID() string                 { return m.id }
func (m *mockCollidable) SetID(id string)            { m.id = id }
func (m *mockCollidable) Position() image.Rectangle  { return image.Rectangle{} }
func (m *mockCollidable) SetPosition(_, _ int)       {}
func (m *mockCollidable) SetPosition16(x, y int)     { m.x16 = x; m.y16 = y }
func (m *mockCollidable) SetSize(_, _ int)           {}
func (m *mockCollidable) Scale() float64             { return 1.0 }
func (m *mockCollidable) SetScale(_ float64)         {}
func (m *mockCollidable) GetPosition16() (int, int)  { return m.x16, m.y16 }
func (m *mockCollidable) GetPositionMin() (int, int) { return 0, 0 }
func (m *mockCollidable) GetShape() body.Shape       { return nil }
func (m *mockCollidable) Owner() interface{}         { return m.owner }
func (m *mockCollidable) SetOwner(o interface{})     { m.owner = o }
func (m *mockCollidable) LastOwner() interface{}     { return nil }
func (m *mockCollidable) OnTouch(_ body.Collidable)  {}
func (m *mockCollidable) OnBlock(_ body.Collidable)  {}
func (m *mockCollidable) GetTouchable() body.Touchable {
	return m
}
func (m *mockCollidable) DrawCollisionBox(_ *ebiten.Image, _ image.Rectangle) {}
func (m *mockCollidable) CollisionPosition() []image.Rectangle                { return nil }
func (m *mockCollidable) CollisionShapes() []body.Collidable                  { return nil }
func (m *mockCollidable) IsObstructive() bool                                 { return false }
func (m *mockCollidable) SetIsObstructive(_ bool)                             {}
func (m *mockCollidable) AddCollision(_ ...body.Collidable)                   {}
func (m *mockCollidable) ClearCollisions()                                    {}
func (m *mockCollidable) SetTouchable(_ body.Touchable)                       {}
func (m *mockCollidable) ApplyValidPosition(_ int, _ bool, _ body.BodiesSpace) (int, int, bool) {
	return 0, 0, false
}

// fakeDamageable implements contractscombat.Damageable and tracks calls.
// It also provides a Faction() method for faction checks.
type fakeDamageable struct {
	takeDamageCalls []int
	faction         combat.Faction
}

func (f *fakeDamageable) TakeDamage(amount int) {
	f.takeDamageCalls = append(f.takeDamageCalls, amount)
}

func (f *fakeDamageable) Faction() combat.Faction {
	return f.faction
}

// fakeDamageableBody implements body.Collidable + contractscombat.Damageable with faction support.
type fakeDamageableBody struct {
	id              string
	owner           interface{}
	takeDamageCalls []int
	faction         combat.Faction
}

func (f *fakeDamageableBody) ID() string                                          { return f.id }
func (f *fakeDamageableBody) SetID(id string)                                     { f.id = id }
func (f *fakeDamageableBody) Position() image.Rectangle                           { return image.Rectangle{} }
func (f *fakeDamageableBody) SetPosition(_, _ int)                                {}
func (f *fakeDamageableBody) SetPosition16(_, _ int)                              {}
func (f *fakeDamageableBody) SetSize(_, _ int)                                    {}
func (f *fakeDamageableBody) Scale() float64                                      { return 1.0 }
func (f *fakeDamageableBody) SetScale(_ float64)                                  {}
func (f *fakeDamageableBody) GetPosition16() (int, int)                           { return 0, 0 }
func (f *fakeDamageableBody) GetPositionMin() (int, int)                          { return 0, 0 }
func (f *fakeDamageableBody) GetShape() body.Shape                                { return nil }
func (f *fakeDamageableBody) Owner() interface{}                                  { return f.owner }
func (f *fakeDamageableBody) SetOwner(o interface{})                              { f.owner = o }
func (f *fakeDamageableBody) LastOwner() interface{}                              { return nil }
func (f *fakeDamageableBody) OnTouch(_ body.Collidable)                           {}
func (f *fakeDamageableBody) OnBlock(_ body.Collidable)                           {}
func (f *fakeDamageableBody) GetTouchable() body.Touchable                        { return f }
func (f *fakeDamageableBody) DrawCollisionBox(_ *ebiten.Image, _ image.Rectangle) {}
func (f *fakeDamageableBody) CollisionPosition() []image.Rectangle                { return nil }
func (f *fakeDamageableBody) CollisionShapes() []body.Collidable                  { return nil }
func (f *fakeDamageableBody) IsObstructive() bool                                 { return false }
func (f *fakeDamageableBody) SetIsObstructive(_ bool)                             {}
func (f *fakeDamageableBody) AddCollision(_ ...body.Collidable)                   {}
func (f *fakeDamageableBody) ClearCollisions()                                    {}
func (f *fakeDamageableBody) SetTouchable(_ body.Touchable)                       {}
func (f *fakeDamageableBody) ApplyValidPosition(_ int, _ bool, _ body.BodiesSpace) (int, int, bool) {
	return 0, 0, false
}
func (f *fakeDamageableBody) TakeDamage(amount int) {
	f.takeDamageCalls = append(f.takeDamageCalls, amount)
}
func (f *fakeDamageableBody) Faction() combat.Faction {
	return f.faction
}

// fakeDestructible implements contractscombat.Destructible (Damageable + IsDestroyed).
// It also provides a Faction() method for faction checks.
type fakeDestructible struct {
	takeDamageCalls []int
	destroyed       bool
	faction         combat.Faction
}

func (f *fakeDestructible) TakeDamage(amount int) {
	f.takeDamageCalls = append(f.takeDamageCalls, amount)
}

func (f *fakeDestructible) IsDestroyed() bool {
	return f.destroyed
}

func (f *fakeDestructible) Faction() combat.Faction {
	return f.faction
}

// fakeCollidableWithOwner creates a mockCollidable with a specified owner.
// Useful for testing resolution of Damageable from body owners.
func fakeCollidableWithOwner(owner interface{}) *mockCollidable {
	return &mockCollidable{
		id:    "fake-body",
		owner: owner,
	}
}

// mockPassthroughBody is a Collidable that also implements body.Passthrough.
// Used to verify projectiles skip passthrough entities entirely.
type mockPassthroughBody struct {
	mockCollidable
}

func (m *mockPassthroughBody) IsPassthrough() bool { return true }

// mockPassthroughOwner is a plain struct that implements body.Passthrough,
// used as an owner of a mockCollidable.
type mockPassthroughOwner struct{}

func (m *mockPassthroughOwner) IsPassthrough() bool { return true }

// Compile-time interface assertions to verify mock types satisfy their contracts.
var (
	_ contractscombat.Damageable   = (*fakeDamageable)(nil)
	_ contractscombat.Damageable   = (*fakeDamageableBody)(nil)
	_ contractscombat.Destructible = (*fakeDestructible)(nil)
)
