package gameitems

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/items"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/engine/systems/sprites"
)

type SignpostItem struct {
	items.BaseItem
}

func NewSignpostItem(x, y int) *SignpostItem {
	frameWidth, frameHeight := 32, 36

	var assets sprites.SpriteAssets
	assets = assets.AddSprite(actors.Idle, "assets/images/item-signpost.png")

	sprites, err := sprites.LoadSprites(assets)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: It should be set in a better place (frameRate)
	frameRate := 10
	rect := physics.NewRect(x, y-(frameHeight/2), frameWidth, frameHeight)
	base := items.NewBaseItem(sprites, frameRate, rect)
	// TODO: Remove this
	base.SetID("SIGN-")
	base.SetCollisionArea(rect)
	base.SetTouchable(base)

	return &SignpostItem{BaseItem: *base}
}
