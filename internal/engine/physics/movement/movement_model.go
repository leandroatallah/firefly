package movement

import (
	"fmt"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
)

type MovementModel interface {
	Update(body body.MovableCollidable, space body.BodiesSpace) error
	SetIsScripted(isScripted bool)
}

// Grounded is implemented by movement models that track an on-ground state.
// Models without grounding semantics (e.g. top-down) may simply not implement it;
// callers should treat the absence of Grounded as "always grounded".
type Grounded interface {
	OnGround() bool
}

// InputBlocker is implemented by movement models that can suppress player input.
type InputBlocker interface {
	IsInputBlocked() bool
}

// GravityController is implemented by movement models that expose a gravity toggle.
type GravityController interface {
	SetGravityEnabled(enabled bool)
}

type MovementModelEnum int

func (m MovementModelEnum) String() string {
	MovementModelMap := map[MovementModelEnum]string{
		TopDown:  "TopDown",
		Platform: "Platform",
		BeatEmUp: "BeatEmUp",
	}
	return MovementModelMap[m]
}

const (
	TopDown MovementModelEnum = iota
	Platform
	BeatEmUp
)

func NewMovementModel(model MovementModelEnum, playerMovementBlocker PlayerMovementBlocker) (MovementModel, error) {
	switch model {
	case TopDown:
		return NewTopDownMovementModel(playerMovementBlocker), nil
	case Platform:
		return NewPlatformMovementModel(playerMovementBlocker), nil
	case BeatEmUp:
		return NewBeatEmUpMovementModel(playerMovementBlocker), nil
	default:
		return nil, fmt.Errorf("unknown movement model type")
	}
}
