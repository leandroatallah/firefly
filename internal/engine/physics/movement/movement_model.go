package movement

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

type MovementModel interface {
	Update(body body.MovableCollidable, space body.BodiesSpace) error
	InputHandler(body body.MovableCollidable, space body.BodiesSpace)
	SetIsScripted(isScripted bool)
}

type MovementModelEnum int

func (m MovementModelEnum) String() string {
	MovementModelMap := map[MovementModelEnum]string{
		TopDown:  "TopDown",
		Platform: "Platform",
	}
	return MovementModelMap[m]
}

const (
	TopDown MovementModelEnum = iota
	Platform
)

func NewMovementModel(model MovementModelEnum, playerMovementBlocker PlayerMovementBlocker) (MovementModel, error) {
	switch model {
	case TopDown:
		return NewTopDownMovementModel(playerMovementBlocker), nil
	case Platform:
		return NewPlatformMovementModel(playerMovementBlocker), nil
	default:
		return nil, fmt.Errorf("unknown movement model type")
	}
}
