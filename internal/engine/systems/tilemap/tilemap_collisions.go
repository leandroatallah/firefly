package tilemap

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/physics"
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
	endpointLayer, err := t.FindLayerByName("Endpoint")
	if err != nil {
		return err
	}
	for _, obj := range endpointLayer.Objects {
		obstacle := t.NewObstacleRect(obj, "Endpoint", false)
		obstacle.SetTouchable(triggerEndpoint)
		space.AddBody(obstacle)
	}

	obstacleLayer, err := t.FindLayerByName("Obstacles")
	if err != nil {
		return err
	}
	for _, obj := range obstacleLayer.Objects {
		obstacle := t.NewObstacleRect(obj, "OBSTACLE", true)
		space.AddBody(obstacle)
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
