package movement

import (
	"fmt"

	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
)

type MovementModel interface {
	Update(body body.MovableCollidable, space body.BodiesSpace) error
	SetIsScripted(isScripted bool)
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
