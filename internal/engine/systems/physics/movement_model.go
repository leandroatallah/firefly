package physics

import (
	"fmt"
)

type MovementModel interface {
	Update(body *PhysicsBody, space *Space) error
	InputHandler(body *PhysicsBody)
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

func NewMovementModel(model MovementModelEnum) (MovementModel, error) {
	switch model {
	case TopDown:
		return NewTopDownMovementModel(), nil
	case Platform:
		return NewPlatformMovementModel(), nil
	default:
		return nil, fmt.Errorf("unknown movement model type")
	}
}
