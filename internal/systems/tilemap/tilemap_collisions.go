package tilemap

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/systems/physics"
)

type LayerNameID int

const (
	ObstaclesLayer LayerNameID = iota
	EnemiesLayer
	ItemsLayer
	PlayerStartLayer
	EndpointLayer
)

var LayerNameMap = map[string]LayerNameID{
	"Obstacles":   ObstaclesLayer,
	"Enemies":     EnemiesLayer,
	"Items":       ItemsLayer,
	"PlayerStart": PlayerStartLayer,
	"Endpoint":    EndpointLayer,
}

func (t *Tilemap) CreateCollisionBodies(space *physics.Space, triggerMap map[LayerNameID]physics.Touchable) error {
	if t == nil {
		return fmt.Errorf("the tilemap was not initialized")
	}

	for _, layer := range t.Layers {
		// CreateCollisionBodies only handles layers of objectgroup type
		if layer.Type != "objectgroup" {
			continue
		}

		switch LayerNameMap[layer.Name] {
		// TODO: Move endpoint to the scene and create as Item Coin
		case EndpointLayer:
			for _, obj := range layer.Objects {
				obstacle := t.NewObstacleRect(obj, false)
				if trigger, ok := triggerMap[EndpointLayer]; ok {
					obstacle.SetTouchable(trigger)
				}
				space.AddBody(obstacle)
			}
		// TODO: Move obstacles to the scene and create as Item Coin
		case ObstaclesLayer:
			for _, obj := range layer.Objects {
				obstacle := t.NewObstacleRect(obj, true)
				space.AddBody(obstacle)
			}
		default:
			continue
		}
	}

	return nil
}

func (t *Tilemap) NewObstacleRect(obj *Obstacle, isObstructive bool) *physics.ObstacleRect {
	mapHeight := t.Height * t.Tileheight
	yOffset := config.ScreenHeight - mapHeight

	rect := physics.NewRect(int(obj.X), int(obj.Y)+yOffset, int(obj.Width), int(obj.Height))
	obstacle := physics.NewObstacleRect(rect).AddCollision()
	obstacle.SetIsObstructive(isObstructive)
	return obstacle
}
