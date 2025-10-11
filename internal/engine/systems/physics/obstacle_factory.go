package physics

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/game/constants"
)

type ObstacleType int

type ObstacleFactory interface {
	Create(obstacleType ObstacleType) (body.Obstacle, error)
}

type DefaultObstacleFactory struct{}

func NewDefaultObstacleFactory() *DefaultObstacleFactory {
	return &DefaultObstacleFactory{}
}

const (
	ObstacleWallTop ObstacleType = iota
	ObstacleWallLeft
	ObstacleWallRight
	ObstacleWallDown
)

func (f *DefaultObstacleFactory) Create(obstableType ObstacleType) (body.Obstacle, error) {
	switch obstableType {
	case ObstacleWallTop:
		return NewWallTop(), nil
	case ObstacleWallLeft:
		return NewWallLeft(), nil
	case ObstacleWallRight:
		return NewWallRight(), nil
	case ObstacleWallDown:
		return NewWallDown(), nil
	default:
		return nil, fmt.Errorf("unknown obstacle type")
	}
}

const wallWidth = 20

func NewWallTop() *ObstacleRect {
	return NewObstacleRect(
		NewRect(0, 0, constants.ScreenWidth, wallWidth),
	).AddCollision(
		NewCollisionArea(
			NewRect(0, 0, constants.ScreenWidth, wallWidth),
		),
	)
}

func NewWallLeft() *ObstacleRect {
	return NewObstacleRect(
		NewRect(0, 0, wallWidth, constants.ScreenHeight),
	).AddCollision()
}

func NewWallRight() *ObstacleRect {
	return NewObstacleRect(
		NewRect(constants.ScreenWidth-wallWidth, 0, wallWidth, constants.ScreenHeight),
	).AddCollision()
}
func NewWallDown() *ObstacleRect {
	return NewObstacleRect(
		NewRect(0, constants.ScreenHeight-wallWidth, constants.ScreenWidth, wallWidth),
	).AddCollision()
}
