// Package-local mocks for Red-Phase tests of story 059-thin-game-phase-scenes.
//
// SPEC.md §3 introduces BuildPlayer[T actors.ActorEntity] together with an
// unexported playerWiring interface (SetInventory / SetMelee / GetCharacter).
// The mocks below provide:
//   - mockPlayerWithWiring: implements actors.ActorEntity AND the wiring
//     interface — used to verify SetInventory / SetMelee call counts.
//   - plainActor: implements actors.ActorEntity but NOT the wiring
//     interface — used to verify BuildPlayer is a no-op for such types.
package kitbuilder

import (
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/vfx"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors"
	"github.com/boilerplate/ebiten-template/internal/engine/entity/actors/movement"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
	physicsmovement "github.com/boilerplate/ebiten-template/internal/engine/physics/movement"
	"github.com/boilerplate/ebiten-template/internal/kit/combat/weapon"
	"github.com/hajimehoshi/ebiten/v2"
)

// --- mockPlayerWithWiring satisfies actors.ActorEntity + playerWiring ------

type mockPlayerWithWiring struct {
	setInventoryCalls int
	lastInventory     interface{}
	setMeleeCalls     int
	lastMelee         *weapon.MeleeWeapon
	lastVFX           vfx.Manager
	character         *actors.Character
}

func newMockPlayerWithWiring() *mockPlayerWithWiring {
	return &mockPlayerWithWiring{character: &actors.Character{}}
}

// Wiring methods consumed by playerWiring.
func (m *mockPlayerWithWiring) SetInventory(inv interface{}) {
	m.setInventoryCalls++
	m.lastInventory = inv
}

func (m *mockPlayerWithWiring) SetMelee(w *weapon.MeleeWeapon, vfxMgr vfx.Manager) {
	m.setMeleeCalls++
	m.lastMelee = w
	m.lastVFX = vfxMgr
}

func (m *mockPlayerWithWiring) GetCharacter() *actors.Character { return m.character }

// --- ActorEntity / body / state-machine boilerplate ------------------------

func (m *mockPlayerWithWiring) ID() string                                     { return "mock" }
func (m *mockPlayerWithWiring) SetID(string)                                   {}
func (m *mockPlayerWithWiring) Position() image.Rectangle                      { return image.Rect(0, 0, 16, 16) }
func (m *mockPlayerWithWiring) SetPosition(int, int)                           {}
func (m *mockPlayerWithWiring) SetPosition16(int, int)                         {}
func (m *mockPlayerWithWiring) SetSize(int, int)                               {}
func (m *mockPlayerWithWiring) Scale() float64                                 { return 1 }
func (m *mockPlayerWithWiring) SetScale(float64)                               {}
func (m *mockPlayerWithWiring) GetPosition16() (int, int)                      { return 0, 0 }
func (m *mockPlayerWithWiring) GetPositionMin() (int, int)                     { return 0, 0 }
func (m *mockPlayerWithWiring) GetShape() body.Shape                           { return bodyphysics.NewRect(0, 0, 16, 16) }
func (m *mockPlayerWithWiring) Altitude() int                                  { return 0 }
func (m *mockPlayerWithWiring) SetAltitude(int)                                {}
func (m *mockPlayerWithWiring) Altitude16() int                                { return 0 }
func (m *mockPlayerWithWiring) SetAltitude16(int)                              {}
func (m *mockPlayerWithWiring) Speed() int                                     { return 0 }
func (m *mockPlayerWithWiring) MaxSpeed() int                                  { return 0 }
func (m *mockPlayerWithWiring) SetSpeed(int) error                             { return nil }
func (m *mockPlayerWithWiring) SetMaxSpeed(int) error                          { return nil }
func (m *mockPlayerWithWiring) MovementModel() physicsmovement.MovementModel   { return nil }
func (m *mockPlayerWithWiring) SetMovementModel(physicsmovement.MovementModel) {}
func (m *mockPlayerWithWiring) Health() int                                    { return 0 }
func (m *mockPlayerWithWiring) MaxHealth() int                                 { return 0 }
func (m *mockPlayerWithWiring) SetHealth(int)                                  {}
func (m *mockPlayerWithWiring) SetMaxHealth(int)                               {}
func (m *mockPlayerWithWiring) LoseHealth(int)                                 {}
func (m *mockPlayerWithWiring) RestoreHealth(int)                              {}
func (m *mockPlayerWithWiring) Update(body.BodiesSpace) error                  { return nil }
func (m *mockPlayerWithWiring) Image() *ebiten.Image                           { return nil }
func (m *mockPlayerWithWiring) ImageOptions() *ebiten.DrawImageOptions         { return nil }
func (m *mockPlayerWithWiring) UpdateImageOptions()                            {}
func (m *mockPlayerWithWiring) State() actors.ActorStateEnum                   { return actors.Idle }
func (m *mockPlayerWithWiring) SetState(actors.ActorState)                     {}
func (m *mockPlayerWithWiring) SetMovementState(_ movement.MovementStateEnum, _ body.MovableCollidable, _ ...movement.MovementStateOption) {
}
func (m *mockPlayerWithWiring) SwitchMovementState(movement.MovementStateEnum) {}
func (m *mockPlayerWithWiring) MovementState() movement.MovementState          { return nil }
func (m *mockPlayerWithWiring) NewState(actors.ActorStateEnum) (actors.ActorState, error) {
	return nil, nil
}
func (m *mockPlayerWithWiring) OnMoveLeft(int)              {}
func (m *mockPlayerWithWiring) OnMoveRight(int)             {}
func (m *mockPlayerWithWiring) BlockMovement()              {}
func (m *mockPlayerWithWiring) UnblockMovement()            {}
func (m *mockPlayerWithWiring) IsMovementBlocked() bool     { return false }
func (m *mockPlayerWithWiring) Hurt(int)                    {}
func (m *mockPlayerWithWiring) Owner() interface{}          { return nil }
func (m *mockPlayerWithWiring) SetOwner(interface{})        {}
func (m *mockPlayerWithWiring) LastOwner() interface{}      { return nil }
func (m *mockPlayerWithWiring) Velocity() (int, int)        { return 0, 0 }
func (m *mockPlayerWithWiring) SetVelocity(int, int)        {}
func (m *mockPlayerWithWiring) Acceleration() (int, int)    { return 0, 0 }
func (m *mockPlayerWithWiring) SetAcceleration(int, int)    {}
func (m *mockPlayerWithWiring) VAltitude16() int            { return 0 }
func (m *mockPlayerWithWiring) SetVAltitude16(int)          {}
func (m *mockPlayerWithWiring) AccelerationAltitude() int   { return 0 }
func (m *mockPlayerWithWiring) SetAccelerationAltitude(int) {}
func (m *mockPlayerWithWiring) Immobile() bool              { return false }
func (m *mockPlayerWithWiring) SetImmobile(bool)            {}
func (m *mockPlayerWithWiring) Freeze() bool                { return false }
func (m *mockPlayerWithWiring) SetFreeze(bool)              {}
func (m *mockPlayerWithWiring) MoveX(int)                   {}
func (m *mockPlayerWithWiring) MoveY(int)                   {}
func (m *mockPlayerWithWiring) OnMoveUpLeft(int)            {}
func (m *mockPlayerWithWiring) OnMoveDownLeft(int)          {}
func (m *mockPlayerWithWiring) OnMoveUpRight(int)           {}
func (m *mockPlayerWithWiring) OnMoveDownRight(int)         {}
func (m *mockPlayerWithWiring) OnMoveUp(int)                {}
func (m *mockPlayerWithWiring) OnMoveDown(int)              {}
func (m *mockPlayerWithWiring) ApplyValidPosition(_ int, _ bool, _ body.BodiesSpace) (int, int, bool) {
	return 0, 0, false
}
func (m *mockPlayerWithWiring) CheckMovementDirectionX()                        {}
func (m *mockPlayerWithWiring) TryJump(int)                                     {}
func (m *mockPlayerWithWiring) SetJumpForceMultiplier(float64)                  {}
func (m *mockPlayerWithWiring) JumpForceMultiplier() float64                    { return 1 }
func (m *mockPlayerWithWiring) SetHorizontalInertia(float64)                    {}
func (m *mockPlayerWithWiring) HorizontalInertia() float64                      { return 1 }
func (m *mockPlayerWithWiring) FaceDirection() animation.FacingDirectionEnum    { return 0 }
func (m *mockPlayerWithWiring) SetFaceDirection(animation.FacingDirectionEnum)  {}
func (m *mockPlayerWithWiring) IsIdle() bool                                    { return true }
func (m *mockPlayerWithWiring) IsWalking() bool                                 { return false }
func (m *mockPlayerWithWiring) IsFalling() bool                                 { return false }
func (m *mockPlayerWithWiring) IsGoingUp() bool                                 { return false }
func (m *mockPlayerWithWiring) SetTouchable(body.Touchable)                     {}
func (m *mockPlayerWithWiring) GetTouchable() body.Touchable                    { return nil }
func (m *mockPlayerWithWiring) OnTouch(body.Collidable)                         {}
func (m *mockPlayerWithWiring) OnBlock(body.Collidable)                         {}
func (m *mockPlayerWithWiring) DrawCollisionBox(*ebiten.Image, image.Rectangle) {}
func (m *mockPlayerWithWiring) CollisionPosition() []image.Rectangle            { return nil }
func (m *mockPlayerWithWiring) CollisionShapes() []body.Collidable              { return nil }
func (m *mockPlayerWithWiring) IsObstructive() bool                             { return false }
func (m *mockPlayerWithWiring) SetIsObstructive(bool)                           {}
func (m *mockPlayerWithWiring) AddCollision(...body.Collidable)                 {}
func (m *mockPlayerWithWiring) ClearCollisions()                                {}
func (m *mockPlayerWithWiring) Invulnerable() bool                              { return false }
func (m *mockPlayerWithWiring) SetInvulnerability(bool)                         {}

// --- plainActor: ActorEntity but does NOT satisfy playerWiring -------------

func newMockPlayerNoWiring() *plainActor { return &plainActor{} }

// plainActor is a separate type that implements actors.ActorEntity but
// intentionally does NOT define SetInventory or SetMelee.
type plainActor struct{}

func (m *plainActor) ID() string                                     { return "plain" }
func (m *plainActor) SetID(string)                                   {}
func (m *plainActor) Position() image.Rectangle                      { return image.Rect(0, 0, 16, 16) }
func (m *plainActor) SetPosition(int, int)                           {}
func (m *plainActor) SetPosition16(int, int)                         {}
func (m *plainActor) SetSize(int, int)                               {}
func (m *plainActor) Scale() float64                                 { return 1 }
func (m *plainActor) SetScale(float64)                               {}
func (m *plainActor) GetPosition16() (int, int)                      { return 0, 0 }
func (m *plainActor) GetPositionMin() (int, int)                     { return 0, 0 }
func (m *plainActor) GetShape() body.Shape                           { return bodyphysics.NewRect(0, 0, 16, 16) }
func (m *plainActor) Altitude() int                                  { return 0 }
func (m *plainActor) SetAltitude(int)                                {}
func (m *plainActor) Altitude16() int                                { return 0 }
func (m *plainActor) SetAltitude16(int)                              {}
func (m *plainActor) Speed() int                                     { return 0 }
func (m *plainActor) MaxSpeed() int                                  { return 0 }
func (m *plainActor) SetSpeed(int) error                             { return nil }
func (m *plainActor) SetMaxSpeed(int) error                          { return nil }
func (m *plainActor) MovementModel() physicsmovement.MovementModel   { return nil }
func (m *plainActor) SetMovementModel(physicsmovement.MovementModel) {}
func (m *plainActor) GetCharacter() *actors.Character                { return nil }
func (m *plainActor) Health() int                                    { return 0 }
func (m *plainActor) MaxHealth() int                                 { return 0 }
func (m *plainActor) SetHealth(int)                                  {}
func (m *plainActor) SetMaxHealth(int)                               {}
func (m *plainActor) LoseHealth(int)                                 {}
func (m *plainActor) RestoreHealth(int)                              {}
func (m *plainActor) Update(body.BodiesSpace) error                  { return nil }
func (m *plainActor) Image() *ebiten.Image                           { return nil }
func (m *plainActor) ImageOptions() *ebiten.DrawImageOptions         { return nil }
func (m *plainActor) UpdateImageOptions()                            {}
func (m *plainActor) State() actors.ActorStateEnum                   { return actors.Idle }
func (m *plainActor) SetState(actors.ActorState)                     {}
func (m *plainActor) SetMovementState(_ movement.MovementStateEnum, _ body.MovableCollidable, _ ...movement.MovementStateOption) {
}
func (m *plainActor) SwitchMovementState(movement.MovementStateEnum) {}
func (m *plainActor) MovementState() movement.MovementState          { return nil }
func (m *plainActor) NewState(actors.ActorStateEnum) (actors.ActorState, error) {
	return nil, nil
}
func (m *plainActor) OnMoveLeft(int)              {}
func (m *plainActor) OnMoveRight(int)             {}
func (m *plainActor) BlockMovement()              {}
func (m *plainActor) UnblockMovement()            {}
func (m *plainActor) IsMovementBlocked() bool     { return false }
func (m *plainActor) Hurt(int)                    {}
func (m *plainActor) Owner() interface{}          { return nil }
func (m *plainActor) SetOwner(interface{})        {}
func (m *plainActor) LastOwner() interface{}      { return nil }
func (m *plainActor) Velocity() (int, int)        { return 0, 0 }
func (m *plainActor) SetVelocity(int, int)        {}
func (m *plainActor) Acceleration() (int, int)    { return 0, 0 }
func (m *plainActor) SetAcceleration(int, int)    {}
func (m *plainActor) VAltitude16() int            { return 0 }
func (m *plainActor) SetVAltitude16(int)          {}
func (m *plainActor) AccelerationAltitude() int   { return 0 }
func (m *plainActor) SetAccelerationAltitude(int) {}
func (m *plainActor) Immobile() bool              { return false }
func (m *plainActor) SetImmobile(bool)            {}
func (m *plainActor) Freeze() bool                { return false }
func (m *plainActor) SetFreeze(bool)              {}
func (m *plainActor) MoveX(int)                   {}
func (m *plainActor) MoveY(int)                   {}
func (m *plainActor) OnMoveUpLeft(int)            {}
func (m *plainActor) OnMoveDownLeft(int)          {}
func (m *plainActor) OnMoveUpRight(int)           {}
func (m *plainActor) OnMoveDownRight(int)         {}
func (m *plainActor) OnMoveUp(int)                {}
func (m *plainActor) OnMoveDown(int)              {}
func (m *plainActor) ApplyValidPosition(_ int, _ bool, _ body.BodiesSpace) (int, int, bool) {
	return 0, 0, false
}
func (m *plainActor) CheckMovementDirectionX()                        {}
func (m *plainActor) TryJump(int)                                     {}
func (m *plainActor) SetJumpForceMultiplier(float64)                  {}
func (m *plainActor) JumpForceMultiplier() float64                    { return 1 }
func (m *plainActor) SetHorizontalInertia(float64)                    {}
func (m *plainActor) HorizontalInertia() float64                      { return 1 }
func (m *plainActor) FaceDirection() animation.FacingDirectionEnum    { return 0 }
func (m *plainActor) SetFaceDirection(animation.FacingDirectionEnum)  {}
func (m *plainActor) IsIdle() bool                                    { return true }
func (m *plainActor) IsWalking() bool                                 { return false }
func (m *plainActor) IsFalling() bool                                 { return false }
func (m *plainActor) IsGoingUp() bool                                 { return false }
func (m *plainActor) SetTouchable(body.Touchable)                     {}
func (m *plainActor) GetTouchable() body.Touchable                    { return nil }
func (m *plainActor) OnTouch(body.Collidable)                         {}
func (m *plainActor) OnBlock(body.Collidable)                         {}
func (m *plainActor) DrawCollisionBox(*ebiten.Image, image.Rectangle) {}
func (m *plainActor) CollisionPosition() []image.Rectangle            { return nil }
func (m *plainActor) CollisionShapes() []body.Collidable              { return nil }
func (m *plainActor) IsObstructive() bool                             { return false }
func (m *plainActor) SetIsObstructive(bool)                           {}
func (m *plainActor) AddCollision(...body.Collidable)                 {}
func (m *plainActor) ClearCollisions()                                {}
func (m *plainActor) Invulnerable() bool                              { return false }
func (m *plainActor) SetInvulnerability(bool)                         {}

// stubInventory is an opaque value used to verify SetInventory's argument
// was the same value passed to BuildPlayer.
type stubInventory struct{ id string }
