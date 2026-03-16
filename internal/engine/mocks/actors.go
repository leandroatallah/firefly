package mocks

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/movement"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
)

// MockActor implements actors.ActorEntity for testing
type MockActor struct {
	Id             string
	Pos            image.Rectangle
	SpeedVal       int
	MaxSpeedVal    int
	HealthVal      int
	MaxHealthVal   int
	MovementMdl    physicsmovement.MovementModel
	MovementSt     movement.MovementStateEnum
	IsScripted     bool
	MoveLeftForce  int
	MoveRightForce int
}

func (m *MockActor) ID() string                                           { return m.Id }
func (m *MockActor) SetID(id string)                                      { m.Id = id }
func (m *MockActor) Position() image.Rectangle                            { return m.Pos }
func (m *MockActor) SetPosition(x, y int)                                 { m.Pos = image.Rect(x, y, x+10, y+10) }
func (m *MockActor) SetSize(w, h int) {
	m.Pos.Max.X = m.Pos.Min.X + w
	m.Pos.Max.Y = m.Pos.Min.Y + h
}
func (m *MockActor) Scale() float64 { return 1.0 }
func (m *MockActor) SetScale(s float64) {}
func (m *MockActor) SetPosition16(x16, y16 int)                           { m.SetPosition(x16/16, y16/16) }
func (m *MockActor) GetPosition16() (int, int)                            { return m.Pos.Min.X * 16, m.Pos.Min.Y * 16 }
func (m *MockActor) GetPositionMin() (int, int)                           { return m.Pos.Min.X, m.Pos.Min.Y }
func (m *MockActor) GetShape() body.Shape                                 { return m }
func (m *MockActor) Width() int                                           { return m.Pos.Dx() }
func (m *MockActor) Height() int                                          { return m.Pos.Dy() }
func (m *MockActor) Speed() int                                           { return m.SpeedVal }
func (m *MockActor) MaxSpeed() int                                        { return m.MaxSpeedVal }
func (m *MockActor) SetSpeed(s int) error                                 { m.SpeedVal = s; return nil }
func (m *MockActor) SetMaxSpeed(s int) error                              { m.MaxSpeedVal = s; return nil }
func (m *MockActor) MovementModel() physicsmovement.MovementModel         { return m.MovementMdl }
func (m *MockActor) SetMovementModel(model physicsmovement.MovementModel) { m.MovementMdl = model }
func (m *MockActor) OnMoveLeft(force int)                                 { m.MoveLeftForce = force }
func (m *MockActor) OnMoveRight(force int)                                { m.MoveRightForce = force }
func (m *MockActor) SetMovementState(state movement.MovementStateEnum, target body.MovableCollidable, options ...movement.MovementStateOption) {
	m.MovementSt = state
}
func (m *MockActor) GetCharacter() *actors.Character {
	return &actors.Character{}
}

// Other interface methods (stubs)
func (m *MockActor) Image() *ebiten.Image                                 { return nil }
func (m *MockActor) ImageOptions() *ebiten.DrawImageOptions               { return nil }
func (m *MockActor) UpdateImageOptions()                                  {}
func (m *MockActor) BlockMovement()                                       {}
func (m *MockActor) UnblockMovement()                                     {}
func (m *MockActor) IsMovementBlocked() bool                              { return false }
func (m *MockActor) State() actors.ActorStateEnum                         { return 0 }
func (m *MockActor) SetState(state actors.ActorState)                     {}
func (m *MockActor) SwitchMovementState(state movement.MovementStateEnum) {}
func (m *MockActor) MovementState() movement.MovementState                { return nil }
func (m *MockActor) NewState(state actors.ActorStateEnum) (actors.ActorState, error) {
	return nil, nil
}
func (m *MockActor) Hurt(damage int)                                     {}
func (m *MockActor) OnDie()                                             {}
func (m *MockActor) OnJump()                                            {}
func (m *MockActor) OnLand()                                            {}
func (m *MockActor) OnFall()                                            {}
func (m *MockActor) SetOnJump(f func(image.Point))                      {}
func (m *MockActor) SetOnFall(f func(image.Point))                      {}
func (m *MockActor) SetOnLand(f func(image.Point))                      {}
func (m *MockActor) SetAppContext(_ any)                                {}
func (m *MockActor) AppContext() any                                    { return nil }
func (m *MockActor) Owner() interface{}                                  { return nil }
func (m *MockActor) SetOwner(interface{})                                { }
func (m *MockActor) LastOwner() interface{}                              { return nil }
func (m *MockActor) Update(space body.BodiesSpace) error                 { return nil }
func (m *MockActor) Health() int                                         { return m.HealthVal }
func (m *MockActor) MaxHealth() int                                      { return m.MaxHealthVal }
func (m *MockActor) SetHealth(h int)                                     { m.HealthVal = h }
func (m *MockActor) SetMaxHealth(h int)                                  { m.MaxHealthVal = h }
func (m *MockActor) LoseHealth(d int)                                    { m.HealthVal -= d }
func (m *MockActor) RestoreHealth(h int)                                 { m.HealthVal += h }
func (m *MockActor) Invulnerable() bool                                  { return false }
func (m *MockActor) SetInvulnerability(v bool)                           {}
func (m *MockActor) GetTouchable() body.Touchable                        { return m }
func (m *MockActor) OnTouch(other body.Collidable)                       {}
func (m *MockActor) OnBlock(other body.Collidable)                       {}
func (m *MockActor) DrawCollisionBox(s *ebiten.Image, p image.Rectangle) {}
func (m *MockActor) CollisionPosition() []image.Rectangle                { return []image.Rectangle{m.Pos} }
func (m *MockActor) CollisionShapes() []body.Collidable                  { return nil }
func (m *MockActor) IsObstructive() bool                                 { return true }
func (m *MockActor) SetIsObstructive(v bool)                             {}
func (m *MockActor) AddCollision(list ...body.Collidable)                {}
func (m *MockActor) ClearCollisions()                                    {}
func (m *MockActor) SetTouchable(t body.Touchable)                       {}
func (m *MockActor) ApplyValidPosition(d int, ax bool, sp body.BodiesSpace) (int, int, bool) {
	return m.Pos.Min.X, m.Pos.Min.Y, false
}
func (m *MockActor) MoveX(d int)                                      { m.MoveRightForce = d }
func (m *MockActor) MoveY(d int)                                      { m.MoveRightForce = d }
func (m *MockActor) OnMoveUpLeft(d int)                               { m.MoveLeftForce = d }
func (m *MockActor) OnMoveDownLeft(d int)                             { m.MoveLeftForce = d }
func (m *MockActor) OnMoveUpRight(d int)                              { m.MoveRightForce = d }
func (m *MockActor) OnMoveDownRight(d int)                            { m.MoveRightForce = d }
func (m *MockActor) OnMoveUp(d int)                                   { m.MoveRightForce = d } // Simplified
func (m *MockActor) OnMoveDown(d int)                                 { m.MoveRightForce = d } // Simplified
func (m *MockActor) Velocity() (int, int)                             { return 0, 0 }
func (m *MockActor) SetVelocity(vx, vy int)                           {}
func (m *MockActor) Acceleration() (ax, ay int)                       { return 0, 0 }
func (m *MockActor) SetAcceleration(ax, ay int)                       {}
func (m *MockActor) Immobile() bool                                   { return false }
func (m *MockActor) SetImmobile(i bool)                               {}
func (m *MockActor) SetFreeze(f bool)                                 {}
func (m *MockActor) Freeze() bool                                     { return false }
func (m *MockActor) FaceDirection() animation.FacingDirectionEnum     { return 0 }
func (m *MockActor) SetFaceDirection(v animation.FacingDirectionEnum) {}
func (m *MockActor) IsIdle() bool                                     { return true }
func (m *MockActor) IsWalking() bool                                  { return false }
func (m *MockActor) IsFalling() bool                                  { return false }
func (m *MockActor) IsGoingUp() bool                                  { return false }
func (m *MockActor) CheckMovementDirectionX()                         {}
func (m *MockActor) TryJump(f int)                                    { m.MoveRightForce = f }
func (m *MockActor) SetJumpForceMultiplier(mu float64)                {}
func (m *MockActor) JumpForceMultiplier() float64                     { return 1.0 }
func (m *MockActor) SetHorizontalInertia(i float64)                   {}
func (m *MockActor) HorizontalInertia() float64                       { return 1.0 }

// MockMovementModel implements physicsmovement.MovementModel
type MockMovementModel struct {
	IsScriptedVal bool
}

func (m *MockMovementModel) Update(b body.MovableCollidable, s body.BodiesSpace) error { return nil }
func (m *MockMovementModel) SetIsScripted(is bool)                                     { m.IsScriptedVal = is }
