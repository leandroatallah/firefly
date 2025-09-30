package movement

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/systems/input"
)

type InputMovementState struct {
	BaseMovementState
}

func NewInputMovementState(base BaseMovementState) *InputMovementState {
	return &InputMovementState{BaseMovementState: base}
}

func (s *InputMovementState) Move() {
	if input.IsSomeKeyPressed(ebiten.KeyA, ebiten.KeyLeft) {
		s.actor.OnMoveLeft(s.actor.Speed())
	}
	if input.IsSomeKeyPressed(ebiten.KeyD, ebiten.KeyRight) {
		s.actor.OnMoveRight(s.actor.Speed())
	}
	if input.IsSomeKeyPressed(ebiten.KeyW, ebiten.KeyUp) {
		s.actor.OnMoveUp(s.actor.Speed())
	}
	if input.IsSomeKeyPressed(ebiten.KeyS, ebiten.KeyDown) {
		s.actor.OnMoveDown(s.actor.Speed())
	}
}
