package gamestates_test

import (
	"image"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/animation"
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

type MockInputSource struct {
	DuckHeldFunc            func() bool
	HasCeilingClearanceFunc func() bool
	HorizontalInputFunc     func() int
	JumpPressedFunc         func() bool
	DashPressedFunc         func() bool
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
