package tilemap

import (
	"fmt"
	"log"

	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	bodyphysics "github.com/leandroatallah/firefly/internal/engine/physics/body"
	"github.com/leandroatallah/firefly/internal/engine/physics/space"
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

func (t *Tilemap) CreateCollisionBodies(space *space.Space, triggerEndpoint body.Touchable) {
	endpointLayer, found := t.FindLayerByName("Endpoint")
	if !found {
		log.Printf("Endpoint layer not found in tilemap")
		return
	}

	for _, obj := range endpointLayer.Objects {
		obstacle := t.NewObstacleRect(obj, "Endpoint", false)
		obstacle.SetTouchable(triggerEndpoint)
		space.AddBody(obstacle)
	}

	obstacleLayer, found := t.FindLayerByName("Obstacles")
	if !found {
		log.Printf("Obstacles layer not found in tilemap")
		return
	}

	for _, obj := range obstacleLayer.Objects {
		obstacle := t.NewObstacleRect(obj, "OBSTACLE", true)
		space.AddBody(obstacle)
	}
}

func (t *Tilemap) NewObstacleRect(obj *Obstacle, prefix string, isObstructive bool) *bodyphysics.ObstacleRect {
	y := int(obj.Y)

	rect := bodyphysics.NewRect(int(obj.X), y, int(obj.Width), int(obj.Height))
	o := bodyphysics.NewObstacleRect(rect)
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
