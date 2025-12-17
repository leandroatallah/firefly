package gameitems

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/items"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/engine/systems/sprites"
)

// TODO: It should be in a sprite related package.
func getSprites(assets map[string]items.AssetData) (sprites.SpriteMap, error) {
	var s sprites.SpriteAssets
	for key, value := range assets {
		var state animation.SpriteState
		switch key {
		case "idle":
			state = items.Idle
		default:
			continue
		}
		s = s.AddSprite(state, value.Path)
	}
	result, err := sprites.LoadSprites(s)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// TODO: SpriteData should be in a sprite related package.
func CreateAnimatedItem(id string, data items.SpriteData) (*items.BaseItem, error) {
	assets, err := getSprites(data.Assets)
	if err != nil {
		return nil, err
	}

	rect := physics.NewRect(data.BodyRect.Rect())
	b := items.NewBaseItem(id, assets, rect)
	b.SetFaceDirection(data.FacingDirection)
	b.SetFrameRate(data.FrameRate)

	return b, nil
}

// TODO: Duplicated
type collisionRectSetter interface {
	AddCollisionRect(state items.ItemStateEnum, rect body.Collidable)
}

func SetItemBodies(item items.Item, data items.SpriteData) error {
	item.SetTouchable(item)

	setter, ok := item.(collisionRectSetter)
	if !ok {
		return fmt.Errorf("item must implement collisionRectSetter")
	}

	// Map collisions from sprite data to handle based on state
	for key, assetData := range data.Assets {
		var state items.ItemStateEnum
		switch key {
		case "idle":
			state = items.Idle
		default:
			continue
		}

		for i, r := range assetData.CollisionRects {
			rect := physics.NewCollidableBodyFromRect(physics.NewRect(r.Rect()))
			rect.SetPosition(r.X, r.Y)
			rect.SetID(fmt.Sprintf("%v_COLLISION_RECT_%d", item.ID(), i))
			setter.AddCollisionRect(state, rect)
		}
	}

	return nil
}

func SetItemStats(item items.Item, data items.StatData) error {
	return nil
}
