package gameobstacles

import (
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
)

const wallWidth = 20

func NewWallTop() *physics.ObstacleRect {
	rect := physics.NewRect(0, 0, config.Get().ScreenWidth, wallWidth)
	o := physics.NewObstacleRect(rect)
	o.SetID("WALL-TOP")
	return o
}

func NewWallLeft() *physics.ObstacleRect {
	rect := physics.NewRect(0, 0, wallWidth, config.Get().ScreenHeight)
	o := physics.NewObstacleRect(rect)
	o.SetID("WALL-LEFT")
	return o
}

func NewWallRight() *physics.ObstacleRect {
	rect := physics.NewRect(config.Get().ScreenWidth-wallWidth, 0, wallWidth, config.Get().ScreenHeight)
	o := physics.NewObstacleRect(rect)
	o.SetID("WALL-RIGHT")
	return o
}
func NewWallDown() *physics.ObstacleRect {
	rect := physics.NewRect(0, config.Get().ScreenHeight-wallWidth, config.Get().ScreenWidth, wallWidth)
	o := physics.NewObstacleRect(rect)
	o.SetID("WALL-DOWN")
	return o
}
