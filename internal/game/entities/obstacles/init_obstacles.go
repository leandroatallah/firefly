package gameobstacles

import (
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
)

const (
	ObstacleWallTop physics.ObstacleType = iota
	ObstacleWallLeft
	ObstacleWallRight
	ObstacleWallDown
)

func InitObstacleMap(ctx *core.AppContext) physics.ObstacleMap {
	obstacleMap := map[physics.ObstacleType]func() body.Obstacle{
		ObstacleWallTop: func() body.Obstacle {
			return NewWallTop()
		},
		ObstacleWallLeft: func() body.Obstacle {
			return NewWallLeft()
		},
		ObstacleWallRight: func() body.Obstacle {
			return NewWallRight()
		},
		ObstacleWallDown: func() body.Obstacle {
			return NewWallDown()
		},
	}
	return obstacleMap
}
