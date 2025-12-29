package gameobstacles

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
)

const (
	ObstacleWallTop bodyphysics.ObstacleType = iota
	ObstacleWallLeft
	ObstacleWallRight
	ObstacleWallDown
)

func InitObstacleMap(ctx *app.AppContext) bodyphysics.ObstacleMap {
	obstacleMap := map[bodyphysics.ObstacleType]func() body.Obstacle{
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
