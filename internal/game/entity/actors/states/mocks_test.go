package gamestates_test

import (
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
	contractsbody "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	"github.com/hajimehoshi/ebiten/v2"
)

type MockBody struct {
	SetSizeFunc       func(w, h int)
	SetVelocityFunc   func(vx, vy int)
	VelocityFunc      func() (int, int)
	HeightFunc        func() int
	PositionFunc      func() image.Rectangle
	GetPosition16Func func() (int, int)
	FaceDirectionFunc func() animation.FacingDirectionEnum
	OwnerFunc         func() interface{}
	AccelerationFunc  func() (int, int)
}

func (m *MockBody) SetSize(w, h int) {
	if m.SetSizeFunc != nil {
		m.SetSizeFunc(w, h)
	}
}

func (m *MockBody) SetVelocity(vx, vy int) {
	if m.SetVelocityFunc != nil {
		m.SetVelocityFunc(vx, vy)
	}
}

func (m *MockBody) Velocity() (int, int) {
	if m.VelocityFunc != nil {
		return m.VelocityFunc()
	}
	return 0, 0
}

func (m *MockBody) Height() int {
	if m.HeightFunc != nil {
		return m.HeightFunc()
	}
	return 0
}

func (m *MockBody) Position() image.Rectangle {
	if m.PositionFunc != nil {
		return m.PositionFunc()
	}
	return image.Rectangle{}
}

func (m *MockBody) GetPosition16() (int, int) {
	if m.GetPosition16Func != nil {
		return m.GetPosition16Func()
	}
	return 0, 0
}

func (m *MockBody) FaceDirection() animation.FacingDirectionEnum {
	if m.FaceDirectionFunc != nil {
		return m.FaceDirectionFunc()
	}
	return animation.FaceDirectionRight
}

func (m *MockBody) Owner() interface{} {
	if m.OwnerFunc != nil {
		return m.OwnerFunc()
	}
	return nil
}

func (m *MockBody) Acceleration() (int, int) {
	if m.AccelerationFunc != nil {
		return m.AccelerationFunc()
	}
	return 0, 0
}

// Stub methods to satisfy MovableCollidable interface
func (m *MockBody) MoveX(distance int)                                              {}
func (m *MockBody) MoveY(distance int)                                              {}
func (m *MockBody) OnMoveLeft(distance int)                                         {}
func (m *MockBody) OnMoveUpLeft(distance int)                                       {}
func (m *MockBody) OnMoveDownLeft(distance int)                                     {}
func (m *MockBody) OnMoveRight(distance int)                                        {}
func (m *MockBody) OnMoveUpRight(distance int)                                      {}
func (m *MockBody) OnMoveDownRight(distance int)                                    {}
func (m *MockBody) OnMoveUp(distance int)                                           {}
func (m *MockBody) OnMoveDown(distance int)                                         {}
func (m *MockBody) SetAcceleration(accX, accY int)                                  {}
func (m *MockBody) SetSpeed(speed int) error                                        { return nil }
func (m *MockBody) SetMaxSpeed(maxSpeed int) error                                  { return nil }
func (m *MockBody) Speed() int                                                      { return 0 }
func (m *MockBody) MaxSpeed() int                                                   { return 0 }
func (m *MockBody) Immobile() bool                                                  { return false }
func (m *MockBody) SetImmobile(immobile bool)                                       {}
func (m *MockBody) SetFreeze(freeze bool)                                           {}
func (m *MockBody) Freeze() bool                                                    { return false }
func (m *MockBody) SetFaceDirection(value animation.FacingDirectionEnum)            {}
func (m *MockBody) IsIdle() bool                                                    { return false }
func (m *MockBody) IsWalking() bool                                                 { return false }
func (m *MockBody) IsFalling() bool                                                 { return false }
func (m *MockBody) IsGoingUp() bool                                                 { return false }
func (m *MockBody) CheckMovementDirectionX()                                        {}
func (m *MockBody) TryJump(force int)                                               {}
func (m *MockBody) SetJumpForceMultiplier(multiplier float64)                       {}
func (m *MockBody) JumpForceMultiplier() float64                                    { return 1.0 }
func (m *MockBody) SetHorizontalInertia(inertia float64)                            {}
func (m *MockBody) HorizontalInertia() float64                                      { return 0 }
func (m *MockBody) GetTouchable() contractsbody.Touchable                           { return nil }
func (m *MockBody) DrawCollisionBox(screen *ebiten.Image, position image.Rectangle) {}
func (m *MockBody) CollisionPosition() []image.Rectangle                            { return nil }
func (m *MockBody) CollisionShapes() []contractsbody.Collidable                     { return nil }
func (m *MockBody) IsObstructive() bool                                             { return false }
func (m *MockBody) SetIsObstructive(value bool)                                     {}
func (m *MockBody) AddCollision(list ...contractsbody.Collidable)                   {}
func (m *MockBody) ClearCollisions()                                                {}
func (m *MockBody) SetPosition(x int, y int)                                        {}
func (m *MockBody) SetPosition16(x16, y16 int)                                      {}
func (m *MockBody) SetTouchable(t contractsbody.Touchable)                          {}
func (m *MockBody) ApplyValidPosition(distance16 int, isXAxis bool, space contractsbody.BodiesSpace) (x, y int, wasBlocked bool) {
	return 0, 0, false
}
func (m *MockBody) OnTouch(other contractsbody.Collidable) {}
func (m *MockBody) OnBlock(other contractsbody.Collidable) {}
func (m *MockBody) ID() string                             { return "" }
func (m *MockBody) SetID(id string)                        {}
func (m *MockBody) GetShape() contractsbody.Shape          { return nil }
func (m *MockBody) SetOwner(owner interface{})             {}
func (m *MockBody) LastOwner() interface{}                 { return nil }
func (m *MockBody) GetPositionMin() (x, y int)             { return 0, 0 }
func (m *MockBody) Scale() float64                         { return 1.0 }
func (m *MockBody) SetScale(scale float64)                 {}

type MockInputSource struct {
	DuckHeldFunc            func() bool
	HasCeilingClearanceFunc func() bool
	HorizontalInputFunc     func() int
	JumpPressedFunc         func() bool
	DashPressedFunc         func() bool
	MeleePressedFunc        func() bool
	AimLockHeldFunc         func() bool
	ShootHeldFunc           func() bool
}

func (m *MockInputSource) DuckHeld() bool {
	if m.DuckHeldFunc != nil {
		return m.DuckHeldFunc()
	}
	return false
}

func (m *MockInputSource) HasCeilingClearance() bool {
	if m.HasCeilingClearanceFunc != nil {
		return m.HasCeilingClearanceFunc()
	}
	return false
}

func (m *MockInputSource) HorizontalInput() int {
	if m.HorizontalInputFunc != nil {
		return m.HorizontalInputFunc()
	}
	return 0
}

func (m *MockInputSource) JumpPressed() bool {
	if m.JumpPressedFunc != nil {
		return m.JumpPressedFunc()
	}
	return false
}

func (m *MockInputSource) DashPressed() bool {
	if m.DashPressedFunc != nil {
		return m.DashPressedFunc()
	}
	return false
}

func (m *MockInputSource) MeleePressed() bool {
	if m.MeleePressedFunc != nil {
		return m.MeleePressedFunc()
	}
	return false
}

func (m *MockInputSource) AimLockHeld() bool {
	if m.AimLockHeldFunc != nil {
		return m.AimLockHeldFunc()
	}
	return false
}

func (m *MockInputSource) ShootHeld() bool {
	if m.ShootHeldFunc != nil {
		return m.ShootHeldFunc()
	}
	return false
}
