package gameobstacles

import (
	"github.com/leandroatallah/firefly/internal/engine/app"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
)

const (
	_ bodyphysics.ObstacleType = iota
)

func InitObstacleMap(ctx *app.AppContext) bodyphysics.ObstacleMap {
	obstacleMap := map[bodyphysics.ObstacleType]func() body.Obstacle{}
	return obstacleMap
}
