package gameobstacles

import (
	"github.com/leandroatallah/firefly/internal/engine/data/config"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
)

const wallWidth = 20

func NewWallTop() *bodyphysics.ObstacleRect {
	rect := bodyphysics.NewRect(0, 0, config.Get().ScreenWidth, wallWidth)
	o := bodyphysics.NewObstacleRect(rect)
	o.SetID("WALL-TOP")
	return o
}

func NewWallLeft() *bodyphysics.ObstacleRect {
	rect := bodyphysics.NewRect(0, 0, wallWidth, config.Get().ScreenHeight)
	o := bodyphysics.NewObstacleRect(rect)
	o.SetID("WALL-LEFT")
	return o
}

func NewWallRight() *bodyphysics.ObstacleRect {
	rect := bodyphysics.NewRect(config.Get().ScreenWidth-wallWidth, 0, wallWidth, config.Get().ScreenHeight)
	o := bodyphysics.NewObstacleRect(rect)
	o.SetID("WALL-RIGHT")
	return o
}
func NewWallDown() *bodyphysics.ObstacleRect {
	rect := bodyphysics.NewRect(0, config.Get().ScreenHeight-wallWidth, config.Get().ScreenWidth, wallWidth)
	o := bodyphysics.NewObstacleRect(rect)
	o.SetID("WALL-DOWN")
	return o
}
