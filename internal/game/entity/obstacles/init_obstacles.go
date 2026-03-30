package gameobstacles

import (
	"github.com/boilerplate/ebiten-template/internal/engine/app"
	"github.com/boilerplate/ebiten-template/internal/engine/contracts/body"
	bodyphysics "github.com/boilerplate/ebiten-template/internal/engine/physics/body"
)

const (
	_ bodyphysics.ObstacleType = iota
)

func InitObstacleMap(ctx *app.AppContext) bodyphysics.ObstacleMap {
	obstacleMap := map[bodyphysics.ObstacleType]func() body.Obstacle{}
	return obstacleMap
}
