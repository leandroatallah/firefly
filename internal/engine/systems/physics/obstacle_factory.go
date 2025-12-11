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

// TODO: This file should not implement concrete objects and should have an example file follwing sequence example.
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
	rect := NewRect(0, 0, config.Get().ScreenWidth, wallWidth)
	o := NewObstacleRect(rect)
	o.SetID("WALL-TOP")
	// o.AddCollisionBodies()
	return o
}

func NewWallLeft() *ObstacleRect {
	rect := NewRect(0, 0, wallWidth, config.Get().ScreenHeight)
	o := NewObstacleRect(rect)
	o.SetID("WALL-LEFT")
	// o.AddCollisionBodies()
	return o
}

func NewWallRight() *ObstacleRect {
	rect := NewRect(config.Get().ScreenWidth-wallWidth, 0, wallWidth, config.Get().ScreenHeight)
	o := NewObstacleRect(rect)
	o.SetID("WALL-RIGHT")
	// o.AddCollisionBodies()
	return o
}
func NewWallDown() *ObstacleRect {
	rect := NewRect(0, config.Get().ScreenHeight-wallWidth, config.Get().ScreenWidth, wallWidth)
	o := NewObstacleRect(rect)
	o.SetID("WALL-DOWN")
	// o.AddCollisionBodies()
	return o
}
