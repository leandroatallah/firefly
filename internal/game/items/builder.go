package gameitems

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/items"
	"github.com/leandroatallah/firefly/internal/engine/schemas"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/engine/systems/sprites"
)

func CreateAnimatedItem(id string, data schemas.SpriteData) (*items.BaseItem, error) {
	stateMap := map[string]animation.SpriteState{
		"idle": items.Idle,
	}
	assets, err := sprites.GetSpritesFromAssets(data.Assets, stateMap)
	if err != nil {
		return nil, err
	}

	rect := physics.NewRect(data.BodyRect.Rect())
	b := items.NewBaseItem(id, assets, rect)
	b.SetFaceDirection(data.FacingDirection)
	b.SetFrameRate(data.FrameRate)

	return b, nil
}

type collisionRectSetter interface {
	AddCollisionRect(state items.ItemStateEnum, rect body.Collidable)
}

func SetItemBodies(item items.Item, data schemas.SpriteData) error {
	setter, ok := item.(collisionRectSetter)
	if !ok {
		return fmt.Errorf("item must implement collisionRectSetter")
	}

	stateMap := map[string]animation.SpriteState{
		"idle": items.Idle,
	}

	idProvider := func(assetKey string, index int) string {
		return fmt.Sprintf("%v_COLLISION_RECT_%s_%d", item.ID(), assetKey, index)
	}

	addCollisionRect := func(state animation.SpriteState, rect body.Collidable) {
		itemState, ok := state.(items.ItemStateEnum)
		if !ok {
			// This should not happen if the stateMap is correct
			return
		}
		setter.AddCollisionRect(itemState, rect)
	}

	physics.SetCollisionBodies(item, data, stateMap, idProvider, addCollisionRect)
	return nil
}

func SetItemStats(item items.Item, data items.StatData) error {
	return nil
}
