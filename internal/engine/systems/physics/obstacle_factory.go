package physics

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
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
		NewRect(0, 0, config.Get().ScreenWidth, wallWidth),
	).AddCollision(
		NewCollisionArea(
			NewRect(0, 0, config.Get().ScreenWidth, wallWidth),
		),
	)
}

func NewWallLeft() *ObstacleRect {
	return NewObstacleRect(
		NewRect(0, 0, wallWidth, config.Get().ScreenHeight),
	).AddCollision()
}

func NewWallRight() *ObstacleRect {
	return NewObstacleRect(
		NewRect(config.Get().ScreenWidth-wallWidth, 0, wallWidth, config.Get().ScreenHeight),
	).AddCollision()
}
func NewWallDown() *ObstacleRect {
	return NewObstacleRect(
		NewRect(0, config.Get().ScreenHeight-wallWidth, config.Get().ScreenWidth, wallWidth),
	).AddCollision()
}
