package skill_test

import (
	"image"
	"testing"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/combat"
	"github.com/boilerplate/ebiten-template/internal/engine/input"
	"github.com/boilerplate/ebiten-template/internal/engine/mocks"
	"github.com/boilerplate/ebiten-template/internal/engine/physics/skill"
	"github.com/hajimehoshi/ebiten/v2"
)

// Red Phase Tests: Inventory-Aware Shooting Skill

func TestShootingSkill_FireDelegatesToActiveWeapon(t *testing.T) {
	fireCalled := false
	mockWeapon := &mocks.MockWeapon{
		CanFireFunc: func() bool { return true },
		FireFunc: func(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection) {
			fireCalled = true
		},
	}

	mockInventory := &mocks.MockInventory{
		ActiveWeaponFunc: func() combat.Weapon { return mockWeapon },
	}

	s := skill.NewShootingSkill(mockInventory)
	body := &mockMovableCollidable{
		getPosition16Func: func() (int, int) { return 100 << 4, 50 << 4 },
		faceDirectionFunc: func() animation.FacingDirectionEnum { return animation.FaceDirectionRight },
	}

	// Inject shoot command
	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{Shoot: true}
	}

	s.HandleInput(body, nil, nil)

	if !fireCalled {
		t.Error("expected weapon.Fire() to be called, but it was not")
	}
}

func TestShootingSkill_NoFireWhenWeaponOnCooldown(t *testing.T) {
	fireCalled := false
	mockWeapon := &mocks.MockWeapon{
		CanFireFunc: func() bool { return false },
		FireFunc: func(x16, y16 int, faceDir animation.FacingDirectionEnum, direction body.ShootDirection) {
			fireCalled = true
		},
	}

	mockInventory := &mocks.MockInventory{
		ActiveWeaponFunc: func() combat.Weapon { return mockWeapon },
	}

	s := skill.NewShootingSkill(mockInventory)
	body := &mockMovableCollidable{
		getPosition16Func: func() (int, int) { return 100 << 4, 50 << 4 },
		faceDirectionFunc: func() animation.FacingDirectionEnum { return animation.FaceDirectionRight },
	}

	// Inject shoot command
	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{Shoot: true}
	}

	s.HandleInput(body, nil, nil)

	if fireCalled {
		t.Error("expected weapon.Fire() NOT to be called when CanFire() is false, but it was")
	}
}

func TestShootingSkill_NoFireWhenInventoryEmpty(t *testing.T) {
	mockInventory := &mocks.MockInventory{
		ActiveWeaponFunc: func() combat.Weapon { return nil },
	}

	s := skill.NewShootingSkill(mockInventory)
	body := &mockMovableCollidable{
		getPosition16Func: func() (int, int) { return 100 << 4, 50 << 4 },
		faceDirectionFunc: func() animation.FacingDirectionEnum { return animation.FaceDirectionRight },
	}

	// Inject shoot command
	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{Shoot: true}
	}

	// Should not panic
	s.HandleInput(body, nil, nil)
}

func TestShootingSkill_WeaponSwitchingOnInput(t *testing.T) {
	switchNextCalled := false
	switchPrevCalled := false

	mockInventory := &mocks.MockInventory{
		ActiveWeaponFunc: func() combat.Weapon { return nil },
		SwitchNextFunc: func() {
			switchNextCalled = true
		},
		SwitchPrevFunc: func() {
			switchPrevCalled = true
		},
	}

	s := skill.NewShootingSkill(mockInventory)
	body := &mockMovableCollidable{
		getPosition16Func: func() (int, int) { return 100 << 4, 50 << 4 },
		faceDirectionFunc: func() animation.FacingDirectionEnum { return animation.FaceDirectionRight },
	}

	// Test WeaponNext
	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{WeaponNext: true}
	}

	s.HandleInput(body, nil, nil)

	if !switchNextCalled {
		t.Error("expected inv.SwitchNext() to be called on WeaponNext input, but it was not")
	}

	// Test WeaponPrev
	switchNextCalled = false
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{WeaponPrev: true}
	}

	s.HandleInput(body, nil, nil)

	if !switchPrevCalled {
		t.Error("expected inv.SwitchPrev() to be called on WeaponPrev input, but it was not")
	}
}

func TestShootingSkill_UpdateHandlesShootRelease(t *testing.T) {
	transitionFromShootingCalled := false
	handler := &mocks.MockStateTransitionHandler{
		TransitionFromShootingFunc: func() {
			transitionFromShootingCalled = true
		},
	}

	mockInventory := &mocks.MockInventory{
		ActiveWeaponFunc: func() combat.Weapon { return nil },
	}

	s := skill.NewShootingSkill(mockInventory)
	s.SetStateTransitionHandler(handler)

	body := &mockMovableCollidable{
		getPosition16Func: func() (int, int) { return 0, 0 },
		faceDirectionFunc: func() animation.FacingDirectionEnum { return animation.FaceDirectionRight },
	}

	oldReader := input.CommandsReader
	defer func() { input.CommandsReader = oldReader }()

	// Simulate shoot held
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{Shoot: true}
	}
	s.Update(body, nil)

	// Simulate shoot released
	input.CommandsReader = func() input.PlayerCommands {
		return input.PlayerCommands{Shoot: false}
	}
	s.Update(body, nil)

	if !transitionFromShootingCalled {
		t.Error("expected TransitionFromShooting() to be called on shoot release, but it was not")
	}
}

type mockMovableCollidable struct {
	getPosition16Func func() (int, int)
	faceDirectionFunc func() animation.FacingDirectionEnum
	isDuckingFunc     func() bool
	width             int
}

type mockShape struct {
	width int
}

func (s *mockShape) Width() int  { return s.width }
func (s *mockShape) Height() int { return 0 }

func (m *mockMovableCollidable) GetPosition16() (int, int) {
	if m.getPosition16Func != nil {
		return m.getPosition16Func()
	}
	return 0, 0
}

func (m *mockMovableCollidable) FaceDirection() animation.FacingDirectionEnum {
	if m.faceDirectionFunc != nil {
		return m.faceDirectionFunc()
	}
	return animation.FaceDirectionRight
}

func (m *mockMovableCollidable) IsDucking() bool {
	if m.isDuckingFunc != nil {
		return m.isDuckingFunc()
	}
	return false
}

func (m *mockMovableCollidable) MoveX(distance int)                              {}
func (m *mockMovableCollidable) MoveY(distance int)                              {}
func (m *mockMovableCollidable) OnMoveLeft(distance int)                         {}
func (m *mockMovableCollidable) OnMoveUpLeft(distance int)                       {}
func (m *mockMovableCollidable) OnMoveDownLeft(distance int)                     {}
func (m *mockMovableCollidable) OnMoveRight(distance int)                        {}
func (m *mockMovableCollidable) OnMoveUpRight(distance int)                      {}
func (m *mockMovableCollidable) OnMoveDownRight(distance int)                    {}
func (m *mockMovableCollidable) OnMoveUp(distance int)                           {}
func (m *mockMovableCollidable) OnMoveDown(distance int)                         {}
func (m *mockMovableCollidable) Velocity() (vx16, vy16 int)                      { return 0, 0 }
func (m *mockMovableCollidable) SetVelocity(vx16, vy16 int)                      {}
func (m *mockMovableCollidable) Acceleration() (accX, accY int)                  { return 0, 0 }
func (m *mockMovableCollidable) SetAcceleration(accX, accY int)                  {}
func (m *mockMovableCollidable) SetSpeed(speed int) error                        { return nil }
func (m *mockMovableCollidable) SetMaxSpeed(maxSpeed int) error                  { return nil }
func (m *mockMovableCollidable) GetPosition() (int, int)                         { return 0, 0 }
func (m *mockMovableCollidable) SetPosition(x, y int)                            {}
func (m *mockMovableCollidable) GetBounds() (x, y, width, height int)            { return 0, 0, 0, 0 }
func (m *mockMovableCollidable) SetBounds(x, y, width, height int)               {}
func (m *mockMovableCollidable) IsCollidingWith(interface{}) bool                { return false }
func (m *mockMovableCollidable) OnCollision(interface{})                         {}
func (m *mockMovableCollidable) GetTouchable() body.Touchable                    { return m }
func (m *mockMovableCollidable) DrawCollisionBox(*ebiten.Image, image.Rectangle) {}
func (m *mockMovableCollidable) CollisionPosition() []image.Rectangle            { return nil }
func (m *mockMovableCollidable) CollisionShapes() []body.Collidable              { return nil }
func (m *mockMovableCollidable) IsObstructive() bool                             { return false }
func (m *mockMovableCollidable) SetIsObstructive(bool)                           {}
func (m *mockMovableCollidable) AddCollision(...body.Collidable)                 {}
func (m *mockMovableCollidable) ClearCollisions()                                {}
func (m *mockMovableCollidable) SetTouchable(body.Touchable)                     {}
func (m *mockMovableCollidable) OnTouch(body.Collidable)                         {}
func (m *mockMovableCollidable) OnBlock(body.Collidable)                         {}
func (m *mockMovableCollidable) ID() string                                      { return "" }
func (m *mockMovableCollidable) SetID(string)                                    {}
func (m *mockMovableCollidable) Position() image.Rectangle                       { return image.Rectangle{} }
func (m *mockMovableCollidable) SetPosition16(x16, y16 int)                      {}
func (m *mockMovableCollidable) SetSize(w, h int)                                {}
func (m *mockMovableCollidable) Scale() float64                                  { return 1.0 }
func (m *mockMovableCollidable) SetScale(float64)                                {}
func (m *mockMovableCollidable) GetPositionMin() (int, int)                      { return 0, 0 }
func (m *mockMovableCollidable) GetShape() body.Shape {
	if m.width > 0 {
		return &mockShape{width: m.width}
	}
	return &mockShape{width: 16} // default width
}
func (m *mockMovableCollidable) ApplyValidPosition(int, bool, body.BodiesSpace) (int, int, bool) {
	return 0, 0, false
}
func (m *mockMovableCollidable) Owner() interface{}                             { return nil }
func (m *mockMovableCollidable) SetOwner(interface{})                           {}
func (m *mockMovableCollidable) LastOwner() interface{}                         { return nil }
func (m *mockMovableCollidable) Speed() int                                     { return 0 }
func (m *mockMovableCollidable) MaxSpeed() int                                  { return 0 }
func (m *mockMovableCollidable) Immobile() bool                                 { return false }
func (m *mockMovableCollidable) SetImmobile(bool)                               {}
func (m *mockMovableCollidable) SetFreeze(bool)                                 {}
func (m *mockMovableCollidable) Freeze() bool                                   { return false }
func (m *mockMovableCollidable) SetFaceDirection(animation.FacingDirectionEnum) {}
func (m *mockMovableCollidable) IsIdle() bool                                   { return false }
func (m *mockMovableCollidable) IsWalking() bool                                { return false }
func (m *mockMovableCollidable) IsFalling() bool                                { return false }
func (m *mockMovableCollidable) IsGoingUp() bool                                { return false }
func (m *mockMovableCollidable) CheckMovementDirectionX()                       {}
func (m *mockMovableCollidable) TryJump(int)                                    {}
func (m *mockMovableCollidable) SetJumpForceMultiplier(float64)                 {}
func (m *mockMovableCollidable) JumpForceMultiplier() float64                   { return 1.0 }
func (m *mockMovableCollidable) SetHorizontalInertia(float64)                   {}
func (m *mockMovableCollidable) HorizontalInertia() float64                     { return 1.0 }
