package mocks

import "github.com/boilerplate/ebiten-template/internal/engine/contracts/body"

// MockStateTransitionHandler is a shared mock for the StateTransitionHandler interface.
type MockStateTransitionHandler struct {
	TransitionToShootingFunc   func(direction body.ShootDirection)
	TransitionFromShootingFunc func()
}

func (m *MockStateTransitionHandler) TransitionToShooting(direction body.ShootDirection) {
	if m.TransitionToShootingFunc != nil {
		m.TransitionToShootingFunc(direction)
	}
}

func (m *MockStateTransitionHandler) TransitionFromShooting() {
	if m.TransitionFromShootingFunc != nil {
		m.TransitionFromShootingFunc()
	}
}
