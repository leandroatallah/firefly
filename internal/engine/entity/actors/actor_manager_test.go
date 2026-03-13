package actors_test

import (
	"image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors"
	"github.com/leandroatallah/firefly/internal/engine/entity/actors/movement"
	physicsmovement "github.com/leandroatallah/firefly/internal/engine/physics/movement"
)

// mockActor implements actors.ActorEntity for testing
type mockActor struct {
	id             string
	pos            image.Rectangle
	speedVal       int
	maxSpeedVal    int
	movementMdl    physicsmovement.MovementModel
	movementSt     movement.MovementStateEnum
	moveLeftForce  int
	moveRightForce int
}

// mockShape implements body.Shape for testing
type mockShape struct {
	width, height int
}

func (m *mockShape) Width() int  { return m.width }
func (m *mockShape) Height() int { return m.height }

func (m *mockActor) ID() string                                           { return m.id }
func (m *mockActor) SetID(id string)                                      { m.id = id }
func (m *mockActor) Position() image.Rectangle                            { return m.pos }
func (m *mockActor) SetPosition(x, y int)                                 { m.pos = image.Rect(x, y, x+10, y+10) }
func (m *mockActor) SetPosition16(x16, y16 int)                           { m.SetPosition(x16/16, y16/16) }
func (m *mockActor) GetPosition16() (int, int)                            { return m.pos.Min.X * 16, m.pos.Min.Y * 16 }
func (m *mockActor) GetPositionMin() (int, int)                           { return m.pos.Min.X, m.pos.Min.Y }
func (m *mockActor) GetShape() body.Shape                                 { return &mockShape{width: m.pos.Dx(), height: m.pos.Dy()} }
func (m *mockActor) Width() int                                           { return m.pos.Dx() }
func (m *mockActor) Height() int                                          { return m.pos.Dy() }
func (m *mockActor) Speed() int                                           { return m.speedVal }
func (m *mockActor) MaxSpeed() int                                        { return m.maxSpeedVal }
func (m *mockActor) SetSpeed(s int) error                                 { m.speedVal = s; return nil }
func (m *mockActor) SetMaxSpeed(s int) error                              { m.maxSpeedVal = s; return nil }
func (m *mockActor) MovementModel() physicsmovement.MovementModel         { return m.movementMdl }
func (m *mockActor) SetMovementModel(model physicsmovement.MovementModel) { m.movementMdl = model }
func (m *mockActor) OnMoveLeft(force int)                                 { m.moveLeftForce = force }
func (m *mockActor) OnMoveRight(force int)                                { m.moveRightForce = force }
func (m *mockActor) SetMovementState(state movement.MovementStateEnum, target body.MovableCollidable, options ...movement.MovementStateOption) {
	m.movementSt = state
}
func (m *mockActor) GetCharacter() *actors.Character {
	return &actors.Character{}
}

// Other interface methods (stubs)
func (m *mockActor) Image() *ebiten.Image                                 { return nil }
func (m *mockActor) ImageOptions() *ebiten.DrawImageOptions               { return nil }
func (m *mockActor) UpdateImageOptions()                                  {}
func (m *mockActor) BlockMovement()                                       {}
func (m *mockActor) UnblockMovement()                                     {}
func (m *mockActor) IsMovementBlocked() bool                              { return false }
func (m *mockActor) State() actors.ActorStateEnum                         { return 0 }
func (m *mockActor) SetState(state actors.ActorState)                     {}
func (m *mockActor) SwitchMovementState(state movement.MovementStateEnum) {}
func (m *mockActor) MovementState() movement.MovementState                { return nil }
func (m *mockActor) NewState(state actors.ActorStateEnum) (actors.ActorState, error) {
	return nil, nil
}
func (m *mockActor) Hurt(damage int)                                     {}
func (m *mockActor) Owner() interface{}                                  { return nil }
func (m *mockActor) SetOwner(interface{})                                { }
func (m *mockActor) LastOwner() interface{}                              { return nil }
func (m *mockActor) Update(space body.BodiesSpace) error                 { return nil }
func (m *mockActor) Health() int                                         { return 100 }
func (m *mockActor) MaxHealth() int                                      { return 100 }
func (m *mockActor) SetHealth(h int)                                     {}
func (m *mockActor) SetMaxHealth(h int)                                  {}
func (m *mockActor) LoseHealth(d int)                                    {}
func (m *mockActor) RestoreHealth(h int)                                 {}
func (m *mockActor) Invulnerable() bool                                  { return false }
func (m *mockActor) SetInvulnerability(v bool)                           {}
func (m *mockActor) GetTouchable() body.Touchable                        { return m }
func (m *mockActor) OnTouch(other body.Collidable)                       {}
func (m *mockActor) OnBlock(other body.Collidable)                       {}
func (m *mockActor) DrawCollisionBox(s *ebiten.Image, p image.Rectangle) {}
func (m *mockActor) CollisionPosition() []image.Rectangle                { return []image.Rectangle{m.pos} }
func (m *mockActor) CollisionShapes() []body.Collidable                  { return nil }
func (m *mockActor) IsObstructive() bool                                 { return true }
func (m *mockActor) SetIsObstructive(v bool)                             {}
func (m *mockActor) AddCollision(list ...body.Collidable)                {}
func (m *mockActor) ClearCollisions()                                    {}
func (m *mockActor) SetTouchable(t body.Touchable)                       {}
func (m *mockActor) ApplyValidPosition(d int, ax bool, sp body.BodiesSpace) (int, int, bool) {
	return m.pos.Min.X, m.pos.Min.Y, false
}
func (m *mockActor) MoveX(d int)                                      { m.moveRightForce = d }
func (m *mockActor) MoveY(d int)                                      { m.moveRightForce = d }
func (m *mockActor) OnMoveUpLeft(d int)                               { m.moveLeftForce = d }
func (m *mockActor) OnMoveDownLeft(d int)                             { m.moveLeftForce = d }
func (m *mockActor) OnMoveUpRight(d int)                              { m.moveRightForce = d }
func (m *mockActor) OnMoveDownRight(d int)                            { m.moveRightForce = d }
func (m *mockActor) OnMoveUp(d int)                                   {}
func (m *mockActor) OnMoveDown(d int)                                 {}
func (m *mockActor) Velocity() (int, int)                             { return 0, 0 }
func (m *mockActor) SetVelocity(vx, vy int)                           {}
func (m *mockActor) Acceleration() (ax, ay int)                       { return 0, 0 }
func (m *mockActor) SetAcceleration(ax, ay int)                       {}
func (m *mockActor) Immobile() bool                                   { return false }
func (m *mockActor) SetImmobile(i bool)                               {}
func (m *mockActor) SetFreeze(f bool)                                 {}
func (m *mockActor) Freeze() bool                                     { return false }
func (m *mockActor) FaceDirection() animation.FacingDirectionEnum     { return 0 }
func (m *mockActor) SetFaceDirection(v animation.FacingDirectionEnum) {}
func (m *mockActor) IsIdle() bool                                     { return true }
func (m *mockActor) IsWalking() bool                                  { return false }
func (m *mockActor) IsFalling() bool                                  { return false }
func (m *mockActor) IsGoingUp() bool                                  { return false }
func (m *mockActor) CheckMovementDirectionX()                         {}
func (m *mockActor) TryJump(f int)                                    { m.moveRightForce = f }
func (m *mockActor) SetJumpForceMultiplier(mu float64)                {}
func (m *mockActor) JumpForceMultiplier() float64                     { return 1.0 }
func (m *mockActor) SetHorizontalInertia(i float64)                   {}
func (m *mockActor) HorizontalInertia() float64                       { return 1.0 }

func TestActorManager(t *testing.T) {
	mgr := actors.NewManager()
	actor1 := &mockActor{id: "a1"}
	actor2 := &mockActor{id: "a2"}
	player := &mockActor{id: "player"}

	mgr.Register(actor1)
	mgr.Register(actor2)
	mgr.Register(player)

	// Test Find
	if a, found := mgr.Find("a1"); !found || a != actor1 {
		t.Error("failed to find actor1")
	}

	// Test GetPlayer
	if p, found := mgr.GetPlayer(); !found || p != player {
		t.Error("failed to find player")
	}

	// Test ForEach
	count := 0
	mgr.ForEach(func(ae actors.ActorEntity) {
		count++
	})
	if count != 3 {
		t.Errorf("expected 3 actors in ForEach, got %d", count)
	}

	// Test Unregister
	mgr.Unregister(actor1)
	if _, found := mgr.Find("a1"); found {
		t.Error("actor1 should have been unregistered")
	}

	// Test Clear
	mgr.Clear()
	if _, found := mgr.GetPlayer(); found {
		t.Error("manager should be empty after Clear")
	}
}
