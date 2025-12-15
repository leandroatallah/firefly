package gameitems

import (
	"log"

	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/contracts/body"
	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/items"
	"github.com/leandroatallah/firefly/internal/engine/systems/physics"
	"github.com/leandroatallah/firefly/internal/engine/systems/sprites"
	gameplayer "github.com/leandroatallah/firefly/internal/game/actors/player"
)

// Concrete
type CollectibleCoinItem struct {
	items.BaseItem
}

func NewCollectibleCoinItem(ctx *core.AppContext, x, y int) *CollectibleCoinItem {
	frameWidth, frameHeight := 16, 16

	var assets sprites.SpriteAssets
	assets = assets.AddSprite(actors.Idle, "assets/images/collectible-coin.png")

	sprites, err := sprites.LoadSprites(assets)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: It should be set in a better place (frameRate)
	frameRate := 10
	rect := physics.NewRect(x, y-(frameHeight/2), frameWidth, frameHeight)
	base := items.NewBaseItem(sprites, frameRate, rect)
	// TODO: Improve this. tmj file has body_id for coins but its complex to bring it here. Maybe it could have a incremental index.
	base.SetID("TEMP")
	base.SetPosition(x, y)
	base.SetCollisionArea(rect)
	base.SetTouchable(base)
	base.SetAppContext(ctx)

	return &CollectibleCoinItem{BaseItem: *base}
}

func (c *CollectibleCoinItem) OnTouch(other body.Collidable) {
	if c.IsRemoved() {
		return
	}

	player, found := c.AppContext().ActorManager.GetPlayer()
	if !found {
		return
	}
	coinCollector, ok := player.(gameplayer.CoinCollector)
	if ok {
		c.SetRemoved(true)
		coinCollector.AddCoinCount(1)
	}
}
