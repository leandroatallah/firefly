package gameitems

import (
	"fmt"

	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/contracts/animation"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/items"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/engine/systems/sprites"
)

// TODO: It should be in a sprite related package.
// TODO: Remove actors package from here
func getSprites(assets map[string]actors.AssetData) (sprites.SpriteMap, error) {
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
// TODO: Remove actors package from here
func CreateAnimatedItem(data actors.SpriteData) (*items.BaseItem, error) {
	assets, err := getSprites(data.Assets)
	if err != nil {
		return nil, err
	}

	rect := physics.NewRect(data.BodyRect.Rect())
	b := items.NewBaseItem(assets, rect)
	b.SetFaceDirection(data.FacingDirection)
	b.SetFrameRate(data.FrameRate)

	return b, nil
}

// TODO: Duplicated
type collisionRectSetter interface {
	AddCollisionRect(state items.ItemStateEnum, rect body.Collidable)
}

// TODO: Remove actors package from here
func SetItemBodies(item items.Item, data actors.SpriteData) error {
	// TODO: Improve this. tmj file has body_id for coins but its complex to bring it here. Maybe it could have a incremental index.
	item.SetID("TEMP")

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

// TODO: Remove actors package from here
func SetItemStats(item items.Item, data actors.StatData) error {
	return nil
}
