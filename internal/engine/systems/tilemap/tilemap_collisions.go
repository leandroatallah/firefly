package tilemap

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
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

func (t *Tilemap) CreateCollisionBodies(space *physics.Space, triggerEndpoint body.Touchable) error {
	if t == nil {
		return fmt.Errorf("the tilemap was not initialized")
	}

	for _, layer := range t.Layers {
		// CreateCollisionBodies only handles layers of objectgroup type
		if !layer.Visible || layer.Type != "objectgroup" {
			continue
		}

		switch LayerNameMap[layer.Name] {
		// TODO: Move endpoint to the scene and create as Item Coin
		case EndpointLayer:
			for _, obj := range layer.Objects {
				obstacle := t.NewObstacleRect(obj, "ENDPOINT", false)
				obstacle.SetTouchable(triggerEndpoint)
				space.AddBody(obstacle)
			}
		// TODO: Move obstacles to the scene and create as Item Coin
		case ObstaclesLayer:
			for _, obj := range layer.Objects {
				obstacle := t.NewObstacleRect(obj, "OBSTACLE", true)
				space.AddBody(obstacle)
			}
		default:
			continue
		}
	}

	return nil
}

func (t *Tilemap) NewObstacleRect(obj *Obstacle, prefix string, isObstructive bool) *physics.ObstacleRect {
	y := int(obj.Y)

	rect := physics.NewRect(int(obj.X), y, int(obj.Width), int(obj.Height))
	o := physics.NewObstacleRect(rect)
	o.SetPosition(int(obj.X), y)
	var id string
	for _, p := range obj.Properties {
		if p.Name == "body_id" {
			id = p.Value
			break
		}
	}
	o.SetID(fmt.Sprintf("%v_%v", prefix, id))
	o.AddCollisionBodies()
	o.SetIsObstructive(isObstructive)
	return o
}
