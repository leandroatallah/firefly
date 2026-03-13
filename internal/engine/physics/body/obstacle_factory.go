package body

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
)

type ObstacleType int
type ObstacleMap map[ObstacleType]func() body.Obstacle

type ObstacleFactory interface {
	Create(obstacleType ObstacleType) (body.Obstacle, error)
}

type DefaultObstacleFactory struct {
	obstacleMap ObstacleMap
}

func NewDefaultObstacleFactory() *DefaultObstacleFactory {
	return &DefaultObstacleFactory{}
}

func (f *DefaultObstacleFactory) Create(obstableType ObstacleType) (body.Obstacle, error) {
	obstacleFunc, ok := f.obstacleMap[obstableType]
	if !ok {
		return nil, fmt.Errorf("unknown obstacle type")
	}
	o := obstacleFunc()
	return o, nil
}
