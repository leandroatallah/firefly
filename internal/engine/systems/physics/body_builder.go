package physics

import (
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/schemas"
)

// SetCollisionBodies processes collision rectangles from sprite data and attaches them to an entity.
// It requires the entity to be touchable, a map of states, a function to provide IDs for the collision bodies,
// and a function to add the collision rectangle to the entity.
func SetCollisionBodies(
	entity body.Collidable,
	data schemas.SpriteData,
	stateMap map[string]animation.SpriteState,
	idProvider func(assetKey string, index int) string,
	addCollisionRect func(state animation.SpriteState, rect body.Collidable),
) {
	entity.SetTouchable(entity)

	for key, assetData := range data.Assets {
		state, ok := stateMap[key]
		if !ok {
			continue
		}

		for i, r := range assetData.CollisionRects {
			rect := NewCollidableBodyFromRect(NewRect(r.Rect()))
			rect.SetPosition(r.X, r.Y)
			rect.SetID(idProvider(key, i))
			addCollisionRect(state, rect)
		}
	}
}
